package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bmoller/cherry-o/game"
)

// keys

type mainKeyMap struct {
	AddPlayer    key.Binding
	Play         key.Binding
	Quit         key.Binding
	RemovePlayer key.Binding
}

func (k mainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k mainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.AddPlayer, k.RemovePlayer, k.Play, k.Quit},
	}
}

var mainKeyBinds = mainKeyMap{
	AddPlayer: key.NewBinding(
		key.WithHelp("a", "Add player"),
		key.WithKeys("a"),
	),
	Play: key.NewBinding(
		key.WithHelp("p", "Play game"),
		key.WithKeys("p"),
	),
	Quit: key.NewBinding(
		key.WithHelp("q", "Quit game"),
		key.WithKeys("q", "ctrl+c"),
	),
	RemovePlayer: key.NewBinding(
		key.WithHelp("r", "Remove player"),
		key.WithKeys("r"),
	),
}

type errorKeyMap struct {
	Dismiss key.Binding
}

func (k errorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Dismiss}
}

func (k errorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Dismiss},
	}
}

var errorKeyBinds = errorKeyMap{
	Dismiss: key.NewBinding(
		key.WithHelp("enter", "Dismiss"),
		key.WithKeys("enter"),
	),
}

type addPlayerKeyMap struct {
	Cancel        key.Binding
	SelectionDown key.Binding
	SelectionUp   key.Binding
	Submit        key.Binding
}

func (k addPlayerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Cancel, k.Submit}
}

func (k addPlayerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SelectionUp, k.SelectionDown, k.Submit, k.Cancel},
	}
}

var addPlayerKeyBinds = addPlayerKeyMap{
	Cancel: key.NewBinding(
		key.WithHelp("esc", "Cancel"),
		key.WithKeys("esc"),
	),
	SelectionDown: key.NewBinding(
		key.WithHelp("↓ / j", "Next color"),
		key.WithKeys("down", "j"),
	),
	SelectionUp: key.NewBinding(
		key.WithHelp("↑ / k", "Previous color"),
		key.WithKeys("up", "k"),
	),
	Submit: key.NewBinding(
		key.WithHelp("enter", "Submit"),
		key.WithKeys("enter"),
	),
}

type removePlayerKeyMap struct {
	Cancel        key.Binding
	SelectionDown key.Binding
	SelectionUp   key.Binding
	Submit        key.Binding
}

func (k removePlayerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Submit}
}

func (k removePlayerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SelectionUp, k.SelectionDown, k.Submit, k.Cancel},
	}
}

var removePlayerKeyBinds = removePlayerKeyMap{
	Cancel: key.NewBinding(
		key.WithHelp("esc", "Cancel"),
		key.WithKeys("esc"),
	),
	SelectionDown: key.NewBinding(
		key.WithHelp("↓ / j", "Next player"),
		key.WithKeys("down", "j"),
	),
	SelectionUp: key.NewBinding(
		key.WithHelp("↑ / k", "Previous player"),
		key.WithKeys("up", "k"),
	),
	Submit: key.NewBinding(
		key.WithHelp("enter", "Submit"),
		key.WithKeys("enter"),
	),
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
	player, ok := item.(game.Player)
	if !ok {
		return
	}

	var renderFunc func(string) string
	switch player.Color() {
	case game.Blue:
		renderFunc = styleBlue.Render
	case game.Green:
		renderFunc = styleGreen.Render
	case game.Red:
		renderFunc = styleRed.Render
	case game.Yellow:
		renderFunc = styleYellow.Render
	}
	styledName := renderFunc(player.Name)
	if index == m.Index() {
		fmt.Fprint(w, " > "+styledName)
	} else {
		fmt.Fprint(w, "   "+styledName)
	}
}
