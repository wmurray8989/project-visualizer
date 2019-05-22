package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
	"github.com/wmurray8989/project-visualizer/config"
	"golang.org/x/mobile/event/key"
)

func main() {
	// Read configuration from disk
	conf := config.Read()

	mw := createMainWindow(&conf)
	wnd := nucular.NewMasterWindow(0, "Project Visualizer", mw.update)
	wnd.SetStyle(style.FromTheme(style.DarkTheme, 2.0))
	wnd.Main()
}

type mainWindow struct {
	conf         *config.Config
	usernameEdit nucular.TextEditor
	passwordEdit nucular.TextEditor
	epicEdit     nucular.TextEditor
	assigneeEdit nucular.TextEditor
	statusEdit   nucular.TextEditor

	cards []card
}

func createMainWindow(conf *config.Config) *mainWindow {
	mw := mainWindow{}
	mw.conf = conf

	mw.usernameEdit.Flags = nucular.EditSelectable
	mw.usernameEdit.Buffer = []rune(conf.Username)

	mw.passwordEdit.Flags = nucular.EditField
	mw.passwordEdit.PasswordChar = '*'
	mw.passwordEdit.Buffer = []rune(conf.Password)

	mw.epicEdit.Flags = nucular.EditSelectable
	mw.epicEdit.Buffer = []rune(conf.Epic)

	mw.assigneeEdit.Flags = nucular.EditSelectable
	mw.assigneeEdit.Buffer = []rune(conf.Assignee)

	mw.statusEdit.Flags = nucular.EditSelectable
	mw.statusEdit.Buffer = []rune(conf.Status)

	// mw.resultEdit.Flags = nucular.EditReadOnly | nucular.EditMultiline | nucular.EditSelectable

	return &mw
}

type card struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields struct {
		Summary string `json:"summary"`
		Epic    string `json:"customfield_10009"`
		Status  struct {
			Name string `json:"name"`
		} `json:"status"`
		Creator struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"creator"`
		Reporter struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"reporter"`
		Assignee struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		IssueLinks []struct {
			ID string `json:"id"`
		} `json:"issuelinks"`
	} `json:"fields"`
}

type issueResultStruct struct {
	StartAt    int    `json:"startAt"`
	MaxResults int    `json:"maxResults"`
	Total      int    `json:"total"`
	Issues     []card `json:"issues"`
}

func getIssues(conf *config.Config) issueResultStruct {
	client := http.Client{
		Timeout: time.Duration(time.Second * 10),
	}

	// Build jql
	jql := ""
	if len(conf.Epic) > 0 {
		jql = jql + "\"epic link\"=\"" + conf.Epic + "\""
	}
	if len(conf.Assignee) > 0 {
		if len(jql) > 0 {
			jql = jql + " AND "
		}
		jql = jql + "\"assignee\"=\"" + conf.Assignee + "\""
	}
	if len(conf.Status) > 0 {
		if len(jql) > 0 {
			jql = jql + " AND "
		}
		jql = jql + "\"status\"=\"" + conf.Status + "\""
	}
	if len(jql) == 0 {
		// Require that some filters be set
		return issueResultStruct{}
	}

	req, err := http.NewRequest("GET", "https://crosschx.atlassian.net/rest/api/3/search", nil)
	if err != nil {
		log.Print(err)
		return issueResultStruct{}
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(conf.Username, conf.Password)
	params := req.URL.Query()
	params.Add("jql", jql)
	req.URL.RawQuery = params.Encode()

	fmt.Println(req.URL.RequestURI())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return issueResultStruct{}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return issueResultStruct{}
	}

	var results = issueResultStruct{}
	json.Unmarshal(body, &results)

	return results
}

func (mw *mainWindow) update(w *nucular.Window) {
	// General event handler
	for _, e := range w.Input().Keyboard.Keys {
		if e.Modifiers.String() != key.ModMeta.String() {
			// Only process keys with meta modifier
			continue
		}
		switch e.Rune {
		case 'q':
			w.Close()
			os.Exit(0)
			return
		case 's':
			mw.conf.Write()
			fmt.Println("Config saved to disk")
		}
	}

	w.Row(25).Dynamic(1)
	// Username Input
	w.Label("Username", "LC")
	mw.usernameEdit.Edit(w)
	mw.conf.Username = string(mw.usernameEdit.Buffer)

	// Password Input
	w.Label("Password", "LC")
	mw.passwordEdit.Edit(w)
	mw.conf.Password = string(mw.passwordEdit.Buffer)

	// Epic Input
	w.Label("Epic", "LC")
	mw.epicEdit.Edit(w)
	mw.conf.Epic = string(mw.epicEdit.Buffer)

	// Assignee Input
	w.Label("Assignee", "LC")
	mw.assigneeEdit.Edit(w)
	mw.conf.Assignee = string(mw.assigneeEdit.Buffer)

	// Status Input
	w.Label("Status", "LC")
	mw.statusEdit.Edit(w)
	mw.conf.Status = string(mw.statusEdit.Buffer)

	// Search Button
	if w.ButtonText("Search") {
		results := getIssues(mw.conf)
		mw.cards = results.Issues
		// for _, issue := range results.Issues {
		// 	mw.resultEdit.Buffer = append(
		// 		mw.resultEdit.Buffer,
		// 		[]rune(issue.Self)...,
		// 	)
		// 	mw.resultEdit.Buffer = append(
		// 		mw.resultEdit.Buffer,
		// 		'\n',
		// 	)
		// }
	}

	// Results Display
	w.Label("Results", "LC")
	for _, card := range mw.cards {
		if w.TreePush(nucular.TreeTab, fmt.Sprintf("%s - %s", card.Key, card.Fields.Summary), false) {
			w.Row(25).Dynamic(2)
			w.Label("Status", "LC")
			w.Label(card.Fields.Status.Name, "LC")

			w.Label("Epic", "LC")
			w.Label(card.Fields.Epic, "LC")

			w.Label("Assignee", "LC")
			w.Label(card.Fields.Assignee.DisplayName, "LC")

			w.Label("Reporter", "LC")
			w.Label(card.Fields.Reporter.DisplayName, "LC")

			w.Label("Creator", "LC")
			w.Label(card.Fields.Creator.DisplayName, "LC")

			w.TreePop()
		}
	}
}
