package windows

import (
	"fmt"
	"os"

	"github.com/aarzilli/nucular"
	"github.com/wmurray8989/project-visualizer/bitbucket"
	"github.com/wmurray8989/project-visualizer/config"
	"golang.org/x/mobile/event/key"
)

// Master defines the struct of the main window
type Master struct {
	conf         *config.Config
	usernameEdit nucular.TextEditor
	passwordEdit nucular.TextEditor
	epicEdit     nucular.TextEditor
	assigneeEdit nucular.TextEditor
	statusEdit   nucular.TextEditor

	cards []bitbucket.Card
}

// NewMaster creates and returns a pointer to a new master
func NewMaster(conf *config.Config) *Master {
	mw := Master{}
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

// Update updates the main window
func (mw *Master) Update(w *nucular.Window) {
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
		results := bitbucket.GetIssues(mw.conf)
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
