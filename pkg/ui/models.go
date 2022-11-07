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

type colorItem struct {
	color game.Color
}

func (c colorItem) FilterValue() string {
	return c.color.String()
}

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
	typedItem, ok := item.(colorItem)
	if !ok {
		return
	}

	if index == m.Index() {
		fmt.Fprint(w, typedItem.color)
	} else {
		fmt.Fprint(w, typedItem.color)
	}
}

// logic for player list

type playerItem struct {
	player *game.Player
}

func (p playerItem) FilterValue() string {
	return ""
}

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
	typedItem, ok := item.(playerItem)
	if !ok {
		return
	}

	if index == m.Index() {
		fmt.Fprint(w, typedItem.player.Name)
	} else {
		fmt.Fprint(w, typedItem.player.Name)
	}
}
