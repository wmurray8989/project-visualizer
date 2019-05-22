package main

import (
	"fmt"
	"os"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
	"golang.org/x/mobile/event/key"
)

func main() {
	mw := CreateMainWindow()
	wnd := nucular.NewMasterWindow(0, "Project Visualizer", mw.update)
	wnd.SetStyle(style.FromTheme(style.DarkTheme, 2.0))
	wnd.Main()
}

type mainWindow struct {
	usernameEdit nucular.TextEditor
	passwordEdit nucular.TextEditor
	resultEdit   nucular.TextEditor
}

func CreateMainWindow() *mainWindow {
	mw := mainWindow{}

	mw.passwordEdit.Flags = nucular.EditField
	mw.passwordEdit.PasswordChar = '*'

	mw.resultEdit.Flags = nucular.EditReadOnly | nucular.EditMultiline

	return &mw
}

func getIssues(mw *mainWindow) {
	username := string(mw.usernameEdit.Buffer)
	apiToken := string(mw.passwordEdit.Buffer)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("API Token: %s\n", apiToken)
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
