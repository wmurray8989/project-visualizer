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
	"golang.org/x/mobile/event/key"
)

func main() {
	mw := createMainWindow()
	wnd := nucular.NewMasterWindow(0, "Project Visualizer", mw.update)
	wnd.SetStyle(style.FromTheme(style.DarkTheme, 2.0))
	wnd.Main()
}

type mainWindow struct {
	usernameEdit nucular.TextEditor
	passwordEdit nucular.TextEditor
	resultEdit   nucular.TextEditor
}

func createMainWindow() *mainWindow {
	mw := mainWindow{}

	mw.passwordEdit.Flags = nucular.EditField
	mw.passwordEdit.PasswordChar = '*'

	mw.resultEdit.Flags = nucular.EditReadOnly | nucular.EditMultiline

	return &mw
}

type issueResultStruct struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
	Issues     []struct {
		ID     string `json:"id"`
		Key    string `json:"key"`
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

func getIssues(mw *mainWindow) {
	username := string(mw.usernameEdit.Buffer)
	apiToken := string(mw.passwordEdit.Buffer)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("API Token: %s\n", apiToken)

	client := http.Client{
		Timeout: time.Duration(time.Second * 10),
	}

	req, err := http.NewRequest("GET", "https://crosschx.atlassian.net/rest/api/3/search", nil)
	if err != nil {
		log.Print(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, apiToken)
	params := req.URL.Query()
	params.Add("jql", "\"epic link\"=\"Automation Triggers POC\"")
	req.URL.RawQuery = params.Encode()

	fmt.Println(req.URL.RequestURI())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var results = issueResultStruct{}
	json.Unmarshal(body, &results)

	fmt.Printf("Results: %+v\n", results)

	// fmt.Println(string(body))
	// json.Unmarshal(body, issueStruct)

}

func (mw *mainWindow) update(w *nucular.Window) {
	// General event handler
	for _, e := range w.Input().Keyboard.Keys {
		switch e.Rune {
		case 'q':
			if e.Modifiers.String() == key.ModMeta.String() {
				w.Close()
				os.Exit(0)
				return
			}
		}
	}

	w.Row(25).Dynamic(1)

	// Username Input
	w.Label("Username", "LC")
	mw.usernameEdit.Edit(w)

	// Password Input
	w.Label("API Token", "LC")
	mw.passwordEdit.Edit(w)

	// Get Issues Button
	if w.ButtonText("Get Issues") {
		getIssues(mw)
	}

	// Results Display
	w.Label("Results", "LC")
	w.Row(200).Dynamic(1)
	mw.resultEdit.Edit(w)

}
