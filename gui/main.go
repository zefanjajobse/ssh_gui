package gui

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *tview.Application

func FillTable(table *tview.Table, hosts []host_info) {
	table.Clear()

	v := reflect.ValueOf(host_info{})
	for i := 0; i < v.NumField(); i++ {
		// fill the first row with information for the columns
		table.SetCell(0, i, &tview.TableCell{ Text: v.Type().Field(i).Name, Color: tcell.ColorYellow, Align: tview.AlignLeft, NotSelectable: true, Expansion: 1})
	}

	for iter, host := range hosts {
		// fill all other rows with the .ssh/config file info
		v := reflect.ValueOf(host)
		for i := 0; i < v.NumField(); i++ {
			table.SetCell(iter + 1, i, &tview.TableCell{ Text: v.Field(i).String(), Color: tcell.ColorWhite, Align: tview.AlignLeft, Expansion: 1})
		}
	}

	// Select the last row in the list if selected row no longer exists
	selected_row, _ := table.GetSelection()
	length := table.GetRowCount()
	if selected_row > length {
		table.Select(length, 0)
	}

	// if only the information bar or less is visible, show info message
	if length <= 1 {
		table.SetCell(1, 0, &tview.TableCell{ Text: "No results", Color: tcell.ColorRed, Align: tview.AlignLeft, NotSelectable: true, Expansion: 1})
	}
}

func Start(a *tview.Application) {
	app = a

	table := tview.NewTable().SetFixed(1, 1).SetSelectable(true, false)

	hosts := getHosts()
	FillTable(table, hosts)
	
	input := tview.NewInputField().SetChangedFunc(func(text string) {
		// On input field changed
		filteredHosts := []host_info{};
		for _, host := range hosts {
			if strings.Contains(host.Name, text) {
				filteredHosts = append(filteredHosts, host)
			}
		}
		FillTable(table, filteredHosts)
	})
	input.SetPlaceholder("Search by name...")

	// On other keyboard input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		selected_row, _ := table.GetSelection()
		length := table.GetRowCount()
		switch event.Key() {
			case tcell.KeyDown:
				table.Select((selected_row + 1 + length) % length, 0)
				return nil;
			case tcell.KeyUp:
				table.Select((selected_row - 1 + length) % length, 0)
				return nil;
			case tcell.KeyEnter:
				app.Suspend(func() {
					// row - 1 since the first row is used for column names
					connect(hosts[selected_row - 1].HostName)
				})
			}
		return event
	})
	
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(input, 1, 0, true).
		AddItem(table, 0, 1, false)
	if err := app.SetRoot(flex, true).Run(); err != nil {
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