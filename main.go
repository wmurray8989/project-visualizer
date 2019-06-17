package main

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
	"github.com/wmurray8989/project-visualizer/config"
	"github.com/wmurray8989/project-visualizer/windows"
)

func main() {
	// Read configuration from disk
	conf := config.Read()

	mw := windows.NewMaster(&conf)
	wnd := nucular.NewMasterWindow(0, "Project Visualizer", mw.Update)
	wnd.SetStyle(style.FromTheme(style.WhiteTheme, 2.0))
	wnd.Main()
}
