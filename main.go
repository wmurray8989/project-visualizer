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
	resultEdit   nucular.TextEditor
}

func createMainWindow(conf *config.Config) *mainWindow {
	mw := mainWindow{}
	mw.conf = conf

	mw.usernameEdit.Buffer = []rune(conf.Username)

	mw.passwordEdit.Flags = nucular.EditField
	mw.passwordEdit.PasswordChar = '*'
	mw.passwordEdit.Buffer = []rune(conf.Password)

	mw.resultEdit.Flags = nucular.EditReadOnly | nucular.EditMultiline | nucular.EditSelectable

	return &mw
}

type issueResultStruct struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
	Issues     []struct {
		ID     string `json:"id"`
		Key    string `json:"key"`
		Self   string `json:"self"`
		Fields struct {
			Status struct {
				Name string `json:"name"`
			} `json:"status"`
			Summary    string `json:"summary"`
			IssueLinks []struct {
				ID string `json:"id"`
			} `json:"issuelinks"`
		} `json:"fields"`
	} `json:"issues"`
}

func getIssues(username, password string) issueResultStruct {
	client := http.Client{
		Timeout: time.Duration(time.Second * 10),
	}

	req, err := http.NewRequest("GET", "https://crosschx.atlassian.net/rest/api/3/search", nil)
	if err != nil {
		log.Print(err)
		return issueResultStruct{}
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	params := req.URL.Query()
	params.Add("jql", "\"epic link\"=\"Automation Triggers POC\"")
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

	fmt.Printf("Results: %+v\n", results)

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
	w.Label("API Token", "LC")
	mw.passwordEdit.Edit(w)
	mw.conf.Password = string(mw.passwordEdit.Buffer)

	// Get Issues Button
	if w.ButtonText("Get Issues") {
		mw.resultEdit.Buffer = []rune("")
		results := getIssues(mw.conf.Username, mw.conf.Password)

		for _, issue := range results.Issues {
			mw.resultEdit.Buffer = append(
				mw.resultEdit.Buffer,
				[]rune(issue.Self)...,
			)
			mw.resultEdit.Buffer = append(
				mw.resultEdit.Buffer,
				'\n',
			)
		}
	}

	// Results Display
	w.Label("Results", "LC")
	w.Row(200).Dynamic(1)
	mw.resultEdit.Edit(w)

}
