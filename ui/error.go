package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errorKeyMap struct {
	Dismiss key.Binding
}

var errorKeyBinds = errorKeyMap{
	Dismiss: key.NewBinding(
		key.WithHelp("esc/enter", "Dismiss"),
		key.WithKeys("esc", "enter"),
	),
}

func (k errorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Dismiss}
}

func (k errorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Dismiss},
	}
}

func updateErrorState(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, errorKeyBinds.Dismiss):
			m.err = nil
			m.state = mainState
		}
	}

	return m, cmd
}

var (
	styleErrorMsg = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(red).
		Foreground(red).
		Padding(1, 2)
)

func viewErrorState(m model) string {
	errorContent := lipgloss.Place(mainPane.GetWidth(), mainPane.GetHeight(),
		lipgloss.Center, lipgloss.Center,
		styleErrorMsg.Render(m.err.Error()),
		lipgloss.WithWhitespaceChars("-"),
		lipgloss.WithWhitespaceForeground(yellow))

	return assembleView(renderPlayers(m, -1), renderHelpContent(m, errorKeyBinds), errorContent)
}
