package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type removePlayerKeyMap struct {
	Cancel         key.Binding
	Select         key.Binding
	NextPlayer     key.Binding
	PreviousPlayer key.Binding
}

func (k removePlayerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Cancel}
}

func (k removePlayerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.PreviousPlayer, k.NextPlayer},
		{k.Select, k.Cancel},
	}
}

var removePlayerKeyBinds = removePlayerKeyMap{
	Cancel: key.NewBinding(
		key.WithHelp("esc", "Cancel"),
		key.WithKeys("esc"),
	),
	NextPlayer: key.NewBinding(
		key.WithHelp("↓", "Next player"),
		key.WithKeys("down"),
	),
	PreviousPlayer: key.NewBinding(
		key.WithHelp("↑", "Previous player"),
		key.WithKeys("up"),
	),
	Select: key.NewBinding(
		key.WithHelp("enter", "Confirm"),
		key.WithKeys("enter"),
	),
}

var selectedRemovalIndex int

func updateRemovePlayerState(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		err error
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, removePlayerKeyBinds.Cancel):
			selectedRemovalIndex = 0
			m.state = mainState
		case key.Matches(msg, removePlayerKeyBinds.Select):
			if m.game, err = m.game.RemovePlayer(m.game.Players()[selectedRemovalIndex].Name); err != nil {
				m.err = err
				m.state = errorState
			} else {
				m.state = mainState
			}
			selectedRemovalIndex = 0
		case key.Matches(msg, removePlayerKeyBinds.NextPlayer):
			selectedRemovalIndex++
			if selectedRemovalIndex == len(m.game.Players()) {
				selectedRemovalIndex--
			}
		case key.Matches(msg, removePlayerKeyBinds.PreviousPlayer):
			selectedRemovalIndex--
			if selectedRemovalIndex < 0 {
				selectedRemovalIndex = 0
			}
		}
	}

	return m, cmd
}

func viewRemovePlayerState(m model) string {
	return assembleView(renderPlayers(m, selectedRemovalIndex), renderHelpContent(m, removePlayerKeyBinds), m.turnView.View())
}
