package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bmoller/cherry-o/pkg/game"
)

type state int

const (
	main state = iota
	removePlayer
	addPlayer
	showError
)

type model struct {
	bindHelp     help.Model
	binds        keyMap
	colorList    list.Model
	currentState state
	err          error
	game         game.Game
	nameInput    textinput.Model
	playerList   list.Model
	turnView     viewport.Model
	winner       game.Player
}

func New() tea.Model {
	colorList := list.New(nil, colorsDelegate{}, 10, 6)
	colorList.Title = "Select a Color"
	for _, function := range []func(bool){
		colorList.SetFilteringEnabled,
		colorList.SetShowFilter,
		colorList.SetShowHelp,
		colorList.SetShowPagination,
		colorList.SetShowStatusBar,
	} {
		function(false)
	}
	colorList.SetStatusBarItemName("color", "colors")

	helpModel := help.New()
	helpModel.ShowAll = true
	helpModel.Width = 30

	playerList := list.New(nil, playersDelegate{}, 30, 7)
	playerList.Title = "Players"
	for _, function := range []func(bool){
		playerList.SetFilteringEnabled,
		playerList.SetShowFilter,
		playerList.SetShowHelp,
		playerList.SetShowPagination,
		playerList.SetShowStatusBar,
	} {
		function(false)
	}
	playerList.SetStatusBarItemName("player", "players")

	return model{
		bindHelp:   helpModel,
		binds:      newKeyMap(),
		colorList:  colorList,
		game:       game.Game{},
		nameInput:  textinput.New(),
		playerList: playerList,
		turnView:   viewport.New(50, 25),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentState {
	case showError:
		return updateShowError(msg, m)
	case addPlayer:
		return updateAddPlayer(msg, m)
	default:
		return updateMain(msg, m)
	}
}

func updateShowError(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.binds.Dismiss):
			m.err = nil
			m.currentState = main
		}
	}

	return m, cmd
}

func updateAddPlayer(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
		err  error
	)

	switch keyMsg, ok := msg.(tea.KeyMsg); {
	case ok && keyMsg.Type == tea.KeyEnter:
		// all done; make the call to add a player
		if m.game, err = m.game.AddPlayer(m.nameInput.Value(), m.colorList.SelectedItem().(game.Color)); err != nil {
			m.currentState = showError
			m.err = err
		} else {
			var players []list.Item
			for _, player := range m.game.Players() {
				players = append(players, player)
			}
			cmds = append(cmds, m.playerList.SetItems(players))
			m.currentState = main
		}
		m.nameInput.Reset()
	case ok && (keyMsg.Type == tea.KeyUp || keyMsg.Type == tea.KeyDown):
		// moving color selection up or down
		m.colorList, cmd = m.colorList.Update(msg)
		cmds = append(cmds, cmd)
	}
	// send everything else to the name field
	m.nameInput, cmd = m.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func updateMain(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		err error
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.binds.AddPlayer):
			if m.game.PlayerCount() == game.MaxPlayers {
				m.err = fmt.Errorf("only %d players can play at a time", game.MaxPlayers)
				m.currentState = showError
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
				m.currentState = addPlayer
			}
		case key.Matches(msg, m.binds.Play):
			if turns, winner, err := m.game.Play(); err != nil {
				m.err = err
				m.currentState = showError
			} else {
				m.turnView.SetContent(renderTurns(turns))
				cmd = viewport.Sync(m.turnView)
				m.winner = winner
			}
		case key.Matches(msg, m.binds.Quit):
			cmd = tea.Quit
		case key.Matches(msg, m.binds.RemovePlayer):
			if m.game.PlayerCount() == 0 {
				m.err = errors.New("no players to remove")
				m.currentState = showError
			} else {
				if m.game, err = m.game.RemovePlayer(m.playerList.SelectedItem().(game.Player).Name); err != nil {
					m.err = err
					m.currentState = showError
				} else {
					var players []list.Item
					for _, player := range m.game.Players() {
						players = append(players, player)
					}
					cmd = m.playerList.SetItems(players)
					m.currentState = main
				}
			}
		}
	default:
		m.turnView, cmd = m.turnView.Update(msg)
	}

	return m, cmd
}

func renderTurns(turns []game.Turn) string {
	// spinnerValues = [7]int{1, 2, 3, 4, -2, -2, -10}
	var (
		format   string
		output   strings.Builder
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

		output.WriteString(renderer(fmt.Sprintf(format, turn.Player.Name)))
	}

	return output.String()
}

func (m model) View() string {
	leftSide := lipgloss.JoinVertical(lipgloss.Left,
		playersPane.Render(m.playerList.View()),
		helpPane.Render(lipgloss.JoinVertical(lipgloss.Center,
			helpTitle,
			m.bindHelp.View(m.binds))),
	)
	var rightSide string
	switch m.currentState {
	case showError:
		rightSide = lipgloss.Place(50, 25,
			lipgloss.Center, lipgloss.Center,
			styleErrorMsg.Render(m.err.Error()),
			lipgloss.WithWhitespaceChars("-"),
			lipgloss.WithWhitespaceForeground(yellow))
	case addPlayer:
		rightSide = lipgloss.JoinVertical(lipgloss.Center,
			"Add a Player",
			m.nameInput.View(),
			m.colorList.View())
	default:
		rightSide = m.turnView.View()
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftSide, rightPane.Render(rightSide),
	)
}
