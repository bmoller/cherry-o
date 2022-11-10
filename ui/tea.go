package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bmoller/cherry-o/game"
)

type appState int

const (
	mainState appState = iota
	errorState
	addPlayerState
	removePlayerState
)

type model struct {
	bindHelp     help.Model
	colorList    list.Model
	currentState appState
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
	// Break glass functionality
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	switch m.currentState {
	case errorState:
		return updateErrorState(msg, m)
	case addPlayerState:
		return updateAddPlayerState(msg, m)
	case removePlayerState:
		return updateRemovePlayerState(msg, m)
	default:
		return updateMainState(msg, m)
	}
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
	switch m.currentState {
	case errorState:
		return viewErrorState(m)
	case addPlayerState:
		return viewAddPlayerState(m)
	case removePlayerState:
		return viewRemovePlayerState(m)
	default:
		return viewMainState(m)
	}

	leftSide := lipgloss.JoinVertical(lipgloss.Left,
		playersPane.Render(m.playerList.View()),
		helpPane.Render(lipgloss.JoinVertical(lipgloss.Center,
			helpTitle,
			m.bindHelp.View(mainKeyBinds))),
	)
	var rightSide string
	switch m.currentState {
	case errorState:
		rightSide = lipgloss.Place(50, 25,
			lipgloss.Center, lipgloss.Center,
			styleErrorMsg.Render(m.err.Error()),
			lipgloss.WithWhitespaceChars("-"),
			lipgloss.WithWhitespaceForeground(yellow))
	case addPlayerState:
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
