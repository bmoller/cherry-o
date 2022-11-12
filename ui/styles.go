package ui

import (
	"github.com/charmbracelet/lipgloss"
)

/*
This file should only be used for common styles; state-specific styles belong in their respective files.
*/

var (
	// Color values

	blue = lipgloss.Color("#268bd2")
	//lint:ignore U1000 defining full color scheme
	cyan    = lipgloss.Color("#2aa198")
	green   = lipgloss.Color("#859900")
	magenta = lipgloss.Color("#d33682")
	//lint:ignore U1000 defining full color scheme
	orange = lipgloss.Color("#cb4b16")
	red    = lipgloss.Color("#dc322f")
	//lint:ignore U1000 defining full color scheme
	violet = lipgloss.Color("#6c71c4")
	yellow = lipgloss.Color("#b58900")

	styleBlue = lipgloss.NewStyle().
			Foreground(blue)

	styleGreen = lipgloss.NewStyle().
			Foreground(green)

	styleRed = lipgloss.NewStyle().
			Foreground(red)

	styleYellow = lipgloss.NewStyle().
			Foreground(yellow)

	// Top-level UI components

	// Large, central pane for displaying the main content of the current state, such as an error or turn list.
	mainPane = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(magenta).
			Height(50).
			Margin(1, 2).
			Width(74)

	// Style for the bottom-left panel of the display, intended to be used to display the list of current players.
	playersPane = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(magenta).
			Height(8).
			Margin(0, 2).
			Padding(1, 2).
			Width(30)

	// Style for the bottom-right panel of the display, intended to be used for displaying keybinds.
	helpPane = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(green).
			Height(8).
			Padding(1, 2).
			Width(40)

	helpTitle = lipgloss.NewStyle().
			Margin(0, 16, 1).
			Underline(true).Render("Help")
)
