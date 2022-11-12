package ui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bmoller/cherry-o/game"
)

type mainKeyMap struct {
	AddPlayer    key.Binding
	Play         key.Binding
	Quit         key.Binding
	RemovePlayer key.Binding
	ScrollDown   key.Binding
	ScrollUp     key.Binding
}

func (k mainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k mainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.AddPlayer, k.ScrollUp, k.ScrollDown},
		{k.Play, k.RemovePlayer, k.Quit},
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
	ScrollDown: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "Scroll down"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "Scroll up"),
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
				m.state = errorState
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
				m.state = addPlayerState
			}
		case key.Matches(msg, mainKeyBinds.Play):
			if turns, winner, err := m.game.Play(); err != nil {
				m.err = err
				m.state = errorState
			} else {
				m.turnView.SetContent(renderTurns(turns))
				m.winner = winner
			}
		case key.Matches(msg, mainKeyBinds.Quit):
			cmd = tea.Quit
		case key.Matches(msg, mainKeyBinds.RemovePlayer):
			if m.game.PlayerCount() == 0 {
				m.err = errors.New("no players to remove")
				m.state = errorState
			} else {
				if m.game, err = m.game.RemovePlayer(m.winner.Name); err != nil {
					m.err = err
					m.state = errorState
				} else {
					m.state = mainState
				}
			}
		case key.Matches(msg, mainKeyBinds.ScrollDown, mainKeyBinds.ScrollUp):
			m.turnView, cmd = m.turnView.Update(msg)
		}
	default:
		m.turnView, cmd = m.turnView.Update(msg)
	}

	return m, cmd
}

func renderTurns(turns []game.Turn) string {
	// spinnerValues = [7]int{1, 2, 3, 4, -2, -2, -10}
	var (
		format string
		// output   strings.Builder
		output   string
		renderer func(string) string
	)

	for _, turn := range turns {
		switch turn.Spin {
		case -10:
			format = "Oh no! %s lost 10 cherries!\n"
		case -2:
			format = "Uh-oh, %s lost 2 cherries.\n"
		case 1:
			format = "%s got another cherry.\n"
		case 2:
			format = "Hey, %s got 2 more cherries!\n"
		case 3:
			format = "Yay, %s got 3 more cherries!\n"
		case 4:
			format = "Hooray, %s got 4 more cherries!\n"
		}

		switch turn.Player.Color() {
		case game.Blue:
			renderer = styleBlue.Render
		case game.Green:
			renderer = styleGreen.Render
		case game.Red:
			renderer = styleRed.Render
		case game.Yellow:
			renderer = styleYellow.Render
		}

		//output.WriteString(renderer(fmt.Sprintf(format, turn.Player.Name)))
		output += renderer(fmt.Sprintf(format, turn.Player.Name)) + "\n"
	}

	//return output.String()
	return output
}

func viewMainState(m model) string {
	return assembleView(renderPlayers(m, m.game.Players(), -1), renderHelpContent(m, mainKeyBinds), m.turnView.View())
}
