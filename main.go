package main

import (
	"ssh_gui/main/gui"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication().EnableMouse(true)

	gui.Start(app)

}
