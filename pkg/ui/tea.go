package ui

import (
	"fmt"

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
	game         *game.Game
	nameInput    textinput.Model
	playerList   list.Model
	resultView   viewport.Model
	winner       *game.Player
}

func New() tea.Model {
	colorList := list.New(nil, playersDelegate{}, 10, 4)
	colorList.Title = "Players"
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

	playerList := list.New(nil, playersDelegate{}, 30, 4)
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
		game:       game.New(),
		nameInput:  textinput.New(),
		playerList: playerList,
		resultView: viewport.New(50, 25),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Make sure the viewport is able to process messages about updating
	var resultViewCmd tea.Cmd
	if m.resultView, resultViewCmd = m.resultView.Update(msg); resultViewCmd != nil {
		cmds = append(cmds, resultViewCmd)
	}

	switch m.currentState {
	case showError:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.binds.Dismiss):
				m.err = nil
				m.currentState = main
			}
		}
	default:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.binds.AddPlayer):
				m.currentState = addPlayer
			case key.Matches(msg, m.binds.Play):
				if transcript, winner, err := m.game.Play(); err != nil {
					m.err = err
					m.currentState = showError
				} else {
					m.winner = winner
					results := make([]string, len(transcript))
					for i, turn := range transcript {
						summary := fmt.Sprintf("%s got %d more cherries", turn.Player().Name, turn.Spin())
						switch turn.Player().Color() {
						case game.Blue:
							results[i] = styleTurnBlue.Render(summary)
						case game.Green:
							results[i] = styleTurnGreen.Render(summary)
						case game.Red:
							results[i] = styleTurnRed.Render(summary)
						case game.Yellow:
							results[i] = styleTurnYellow.Render(summary)
						default:
							results[i] = summary
						}
					}
					m.resultView.SetContent(lipgloss.JoinVertical(lipgloss.Left, results...))
					if cmd := viewport.Sync(m.resultView); cmd != nil {
						cmds = append(cmds, cmd)
					}
				}
			case key.Matches(msg, m.binds.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.binds.RemovePlayer):
				player := m.playerList.SelectedItem().(playerItem).player
				if err := m.game.RemovePlayer(player); err != nil {
					m.currentState = showError
					m.err = err
				} else {
					var players []list.Item
					for _, player := range m.game.Players() {
						players = append(players, playerItem{
							player: player,
						})
					}
					if cmd := m.playerList.SetItems(players); cmd != nil {
						cmds = append(cmds, cmd)
					}
					m.currentState = main
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	leftSide := lipgloss.JoinVertical(lipgloss.Left,
		playersPane.Render(m.playerList.View()),
		helpPane.Render(lipgloss.JoinVertical(lipgloss.Center,
			"Help",
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
	default:
		rightSide = m.resultView.View()
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftSide, rightPane.Render(rightSide),
	)
}
