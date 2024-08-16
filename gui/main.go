package gui

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var exited bool
var app *tview.Application

func Start(a *tview.Application) {
	app = a
	table := tview.NewTable().SetFixed(1, 1).SetSelectable(true, false)

	v := reflect.ValueOf(host_info{})
	for i := 0; i < v.NumField(); i++ {
		// fill the first row with information for the columns
		table.SetCell(0, i, &tview.TableCell{ Text: v.Type().Field(i).Name, Color: tcell.ColorYellow, Align: tview.AlignLeft, NotSelectable: true})
	}
	hosts := getHosts()
	for iter, host := range hosts {
		// fill all other rows with the .ssh/config file info
		v := reflect.ValueOf(host)
		for i := 0; i < v.NumField(); i++ {
			table.SetCell(iter + 1, i, &tview.TableCell{ Text: v.Field(i).String(), Color: tcell.ColorWhite, Align: tview.AlignLeft})
		}
	}

	table.SetSelectedFunc(func (row int, column int) {
		// row - 1 since the first row is used for column names
		app.Suspend(func() {
			connect(hosts[row - 1].HostName)
		})
	})

	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

func connect(hostname string) {
	fmt.Println("» ssh", hostname)

	// start ssh session
	cmd := exec.Command("ssh", hostname)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}