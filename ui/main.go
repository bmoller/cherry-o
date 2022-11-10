package ui

import (
	"errors"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bmoller/cherry-o/game"
)

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
		key.WithKeys("a"),
		key.WithHelp("a", "Add player"),
	),
	Play: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Play game"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "Quit game"),
	),
	RemovePlayer: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Remove player"),
	),
}

func updateMainState(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		err error
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, mainKeyBinds.AddPlayer):
			if m.game.PlayerCount() == game.MaxPlayers {
				m.err = fmt.Errorf("only %d players can play at a time", game.MaxPlayers)
				m.currentState = errorState
			} else {
				var (
					cmds      [2]tea.Cmd
					colorList []list.Item
				)
				for color, available := range m.game.AvailableColors() {
					if available {
						colorList = append(colorList, color)
					}
				}
				cmds[0] = m.nameInput.Focus()
				cmds[1] = m.colorList.SetItems(colorList)
				m.colorList.ResetSelected()
				cmd = tea.Batch(cmds[:]...)
				m.currentState = addPlayerState
			}
		case key.Matches(msg, mainKeyBinds.Play):
			if turns, winner, err := m.game.Play(); err != nil {
				m.err = err
				m.currentState = errorState
			} else {
				m.turnView.SetContent(renderTurns(turns))
				cmd = viewport.Sync(m.turnView)
				m.winner = winner
			}
		case key.Matches(msg, mainKeyBinds.Quit):
			cmd = tea.Quit
		case key.Matches(msg, mainKeyBinds.RemovePlayer):
			if m.game.PlayerCount() == 0 {
				m.err = errors.New("no players to remove")
				m.currentState = errorState
			} else {
				if m.game, err = m.game.RemovePlayer(m.playerList.SelectedItem().(game.Player).Name); err != nil {
					m.err = err
					m.currentState = errorState
				} else {
					var players []list.Item
					for _, player := range m.game.Players() {
						players = append(players, player)
					}
					cmd = m.playerList.SetItems(players)
					m.currentState = mainState
				}
			}
		}
	default:
		m.turnView, cmd = m.turnView.Update(msg)
	}

	return m, cmd
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

func viewMainState(m model) string {
	return ""
}
