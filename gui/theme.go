package gui

import (
	"os/exec"
	"runtime"

	"github.com/gdamore/tcell/v2"
)

func checkDarkMode() bool {
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false
		}
	}
	return true
}

type theme struct {
	header            tcell.Color
	app_bg            tcell.Color
	cell              tcell.Color
	cell_bg           tcell.Color
	selected_cell_bg  tcell.Color
	selected_cell     tcell.Color
	input             tcell.Color
	input_bg          tcell.Color
	input_placeholder tcell.Color
}

func GetColorTheme() theme {
	colors := theme{
		header:            tcell.ColorYellow,
		cell_bg:           tcell.ColorBlack,
		cell:              tcell.ColorWhite,
		selected_cell_bg:  tcell.ColorWhite,
		selected_cell:     tcell.ColorBlack,
		app_bg:            tcell.ColorBlack,
		input_bg:          tcell.ColorGray,
		input:             tcell.ColorWhite,
		input_placeholder: tcell.ColorLightGray,
	}
	// use light mode if on mac os and it has dark mode disabled
	if runtime.GOOS == "darwin" && !checkDarkMode() {
		colors = theme{
			header:            tcell.ColorOrange,
			cell_bg:           tcell.ColorWhite,
			cell:              tcell.ColorBlack,
			selected_cell_bg:  tcell.ColorLightGray,
			selected_cell:     tcell.ColorBlack,
			app_bg:            tcell.ColorWhite,
			input_bg:          tcell.ColorLightGray,
			input:             tcell.ColorBlack,
			input_placeholder: tcell.ColorDarkGray,
		}
	}
	return colors
}
