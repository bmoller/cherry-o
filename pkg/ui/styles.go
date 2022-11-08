package ui

import (
	"github.com/charmbracelet/lipgloss"
)

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

	// Top-level UI components

	helpPane = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(green).
			Margin(0, 2).
			Width(30)

	playersPane = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(magenta).
			Height(8).
			Margin(1, 2).
			Padding(1, 2).
			Width(30)

	rightPane = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(magenta).
			Height(25).
			Margin(1, 0).
			Width(50)

	styleErrorMsg = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(red).
			Foreground(red).
			Padding(1, 2)

	styleBlue = lipgloss.NewStyle().
			Foreground(blue)

	styleGreen = lipgloss.NewStyle().
			Foreground(green)

	styleRed = lipgloss.NewStyle().
			Foreground(red)

	styleYellow = lipgloss.NewStyle().
			Foreground(yellow)
)
