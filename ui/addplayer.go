package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bmoller/cherry-o/game"
)

type addPlayerKeyMap struct {
	Cancel        key.Binding
	NextColor     key.Binding
	PreviousColor key.Binding
	Submit        key.Binding
}

func (k addPlayerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Cancel, k.Submit}
}

func (k addPlayerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.PreviousColor, k.NextColor},
		{k.Submit, k.Cancel},
	}
}

var addPlayerKeyBinds = addPlayerKeyMap{
	Cancel: key.NewBinding(
		key.WithHelp("esc", "Cancel"),
		key.WithKeys("esc"),
	),
	NextColor: key.NewBinding(
		key.WithHelp("↓", "Next color"),
		key.WithKeys("down"),
	),
	PreviousColor: key.NewBinding(
		key.WithHelp("↑", "Previous color"),
		key.WithKeys("up"),
	),
	Submit: key.NewBinding(
		key.WithHelp("enter", "Submit"),
		key.WithKeys("enter"),
	),
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
	color, ok := item.(game.Color)
	if !ok {
		return
	}

	var renderFunc func(string) string
	switch color {
	case game.Blue:
		renderFunc = styleBlue.Render
	case game.Green:
		renderFunc = styleGreen.Render
	case game.Red:
		renderFunc = styleRed.Render
	case game.Yellow:
		renderFunc = styleYellow.Render
	}
	itemText := renderFunc(color.String())
	if index == m.Index() {
		fmt.Fprint(w, " > ", itemText)
	} else {
		fmt.Fprint(w, "   ", itemText)
	}
}

func updateAddPlayerState(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		err error
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, addPlayerKeyBinds.Cancel):
			m.nameInput.Reset()
			m.state = mainState
		case key.Matches(msg, addPlayerKeyBinds.Submit):
			// all done; make the call to add a player
			if m.game, err = m.game.AddPlayer(m.nameInput.Value(), m.colorList.SelectedItem().(game.Color)); err != nil {
				m.state = errorState
				m.err = err
			} else {
				var players []list.Item
				for _, player := range m.game.Players() {
					players = append(players, player)
				}
				cmd = m.playerList.SetItems(players)
				m.state = mainState
			}
			m.nameInput.Reset()
		case key.Matches(msg, addPlayerKeyBinds.NextColor, addPlayerKeyBinds.PreviousColor):
			// moving color selection up or down
			m.colorList, cmd = m.colorList.Update(msg)
		default:
			// send everything else to the name field
			m.nameInput, cmd = m.nameInput.Update(msg)
		}
	}

	return m, cmd
}

func viewAddPlayerState(m model) string {
	addPlayerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		"AddPlayer",
		m.nameInput.View(),
		m.colorList.View(),
	)

	return assembleView(m.playerList.View(), renderHelpContent(m, addPlayerKeyBinds), addPlayerContent)
}
