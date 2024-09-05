package gui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *tview.Application

func FillTable(table *tview.Table, hosts []host_info, colors theme) {
	table.Clear()

	cell := &tview.TableCell{ 
		Align: tview.AlignLeft,
		Expansion: 1,
		SelectedStyle: tcell.StyleDefault.Foreground(colors.selected_cell).Background(colors.selected_cell_bg),
		BackgroundColor: colors.cell_bg,
	}

	v := reflect.ValueOf(host_info{})
	for i := 0; i < v.NumField(); i++ {
		current := *cell
		current.Text = v.Type().Field(i).Name
		current.Color = colors.header
		current.NotSelectable = true
		// fill the first row with information for the columns
		table.SetCell(0, i, &current)
	}

	for iter, host := range hosts {
		// fill all other rows with the .ssh/config file info
		v := reflect.ValueOf(host)
		for i := 0; i < v.NumField(); i++ {
			current := *cell
			current.Text = v.Field(i).String()
			current.Color = colors.cell
			current.Clicked = func() bool {
				selected_row, _ := table.GetSelection()
				// row - 1 since the first row is used for column names
				if iter == selected_row - 1 {
					if hosts[iter].Name == "Exit/Quit" {
						app.Stop()
						return false
					}

					app.Suspend(func() {
						connect(hosts[iter])
					})
				}
				return false
			}
			table.SetCell(iter + 1, i, &current)
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
		current := *cell
		current.Text = "No results"
		current.Color = tcell.ColorRed
		current.NotSelectable = true
		table.SetCell(1, 0, &current)
	}
}

func Start(a *tview.Application) {
	app = a
	colors := GetColorTheme()

	table := tview.NewTable().SetFixed(1, 1).SetSelectable(true, false)
	table.SetBackgroundColor(colors.app_bg)
	hosts := getHosts()
	// Add extra option to exit/quit the app
	hosts = append(hosts, host_info{Name: "Exit/Quit"})
	FillTable(table, hosts, colors)
	filteredHosts := hosts;
	
	input := tview.NewInputField().SetChangedFunc(func(text string) {
		// On input field changed
		filteredHosts = []host_info{};
		for _, host := range hosts {
			if strings.Contains(strings.ToLower(host.Name), strings.ToLower(text)) {
				filteredHosts = append(filteredHosts, host)
			}
		}
		FillTable(table, filteredHosts, colors)
		// Select the first row after input field change
		table.Select(0, 0)
	})
	input = input.SetPlaceholder("Search by name...").SetFieldBackgroundColor(colors.input_bg)
	input.SetPlaceholderStyle(input.GetFieldStyle()).SetPlaceholderTextColor(colors.input_placeholder).SetFieldTextColor(colors.input)

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
				// Skip if no row is selected
				if selected_row > len(filteredHosts) || selected_row < 0 {
					return nil;
				}
				app.Suspend(func() {
					// Exit is a extra row after the list of hosts
					if filteredHosts[selected_row - 1].Name == "Exit/Quit" {
						app.Stop()
						return;
					}
					// row - 1 since the first row is used for column names
					connect(filteredHosts[selected_row - 1])
				})
			}
		app.SetFocus(input)
		return event
	})
	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		selected_row, _ := table.GetSelection()
		length := table.GetRowCount()
		switch action {
			case tview.MouseScrollDown:
				table.Select((selected_row + 1 + length) % length, 0)
			case tview.MouseScrollUp:
				table.Select((selected_row - 1 + length) % length, 0)
		}
		app.SetFocus(input)
		return event, action
	})
	
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(input, 1, 0, true).
		AddItem(table, 0, 1, false)
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func connect(host host_info) {
	if host.HostName == "" {
		log.Fatal("Can't connect to specified host, There is no hostname set for ", host.Name)
	}

	cmd := exec.Command("ssh", host.HostName)
	if host.User != "" {
		cmd = exec.Command("ssh", fmt.Sprintf("%s@%s", host.User, host.HostName))
	}
	if host.Port != "22" {
		cmd.Args = append(cmd.Args, "-p", host.Port)
	}
	fmt.Println("» ssh", strings.Join(cmd.Args, " "))
	// start ssh session
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}