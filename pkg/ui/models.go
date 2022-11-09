package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bmoller/cherry-o/pkg/game"
)

type keyMap struct {
	AddPlayer    key.Binding
	Dismiss      key.Binding
	Play         key.Binding
	Quit         key.Binding
	RemovePlayer key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.AddPlayer, k.RemovePlayer, k.Play, k.Dismiss, k.Quit},
	}
}

func newKeyMap() keyMap {
	return keyMap{
		AddPlayer: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Add player"),
		),
		Dismiss: key.NewBinding(
			key.WithKeys("esc", "enter"),
			key.WithHelp("esc/enter", "Dismiss error"),
		),
		Play: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Play game"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "Quit"),
		),
		RemovePlayer: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "Remove player"),
		),
	}
}

// Logic for color list

type colorsDelegate struct{}

func (c colorsDelegate) Height() int {
	return 1
}

func (c colorsDelegate) Spacing() int {
	return 0
}

func (c colorsDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (c colorsDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	typedItem, ok := item.(game.Color)
	if !ok {
		return
	}

	var renderFunc func(string) string
	switch typedItem {
	case game.Blue:
		renderFunc = styleBlue.Render
	case game.Green:
		renderFunc = styleGreen.Render
	case game.Red:
		renderFunc = styleRed.Render
	case game.Yellow:
		renderFunc = styleYellow.Render
	}
	itemText := renderFunc(typedItem.String())
	if index == m.Index() {
		fmt.Fprint(w, " > ", itemText)
	} else {
		fmt.Fprint(w, "   ", itemText)
	}
}

// logic for player list

type playersDelegate struct{}

func (p playersDelegate) Height() int {
	return 1
}

func (p playersDelegate) Spacing() int {
	return 0
}

func (p playersDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (p playersDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	typedItem, ok := item.(game.Player)
	if !ok {
		return
	}

	if index == m.Index() {
		fmt.Fprint(w, typedItem.Name)
	} else {
		fmt.Fprint(w, typedItem.Name)
	}
}
