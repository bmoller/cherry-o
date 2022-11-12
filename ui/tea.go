package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bmoller/cherry-o/game"
)

/*
appState is used to track what the user is currently doing and what to display.

If/when new states and functionality are added, the state should be added here with references to its handler functions.
Each appState needs to have a function for update and view operations, corresponding to the Elm Architecture concepts.
*/
type appState int

const (
	// mainState represents the default view of the application; from here the user can move to any of the other appStates.
	// Defined first so that it acts as default and a sane zero value.
	mainState appState = iota
	// errorState indicates that an error has occurred and displays it to the user.
	errorState
	// addPlayerState shows components used to add a new player to the game.
	addPlayerState
	// removePlayerState handles removal of an existing player from the game, assuming there are any.
	removePlayerState
)

/*
update calls the appState's function for handling event messages.
Because the function receives a state and not the model itself, m must be passed along with msg.
*/
func (s appState) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch s {
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

/*
view calls the appState's function for rendering the current state of the application.
As with update, the model must be passed in m so that components are accessible.
*/
func (s appState) view(m model) string {
	switch s {
	case errorState:
		return viewErrorState(m)
	case addPlayerState:
		return viewAddPlayerState(m)
	case removePlayerState:
		return viewRemovePlayerState(m)
	default:
		return viewMainState(m)
	}
}

type model struct {
	// Used to display the current state's keybinds.
	bindHelp help.Model
	// Presents available colors to the user when adding a new player.
	colorList list.Model
	// The current error resulting in an errorState, if any.
	err error
	// An embedded game simulation; its outputs are presented to the user via the model.
	game game.Game
	// Used to query the name when adding a new player.
	nameInput textinput.Model
	// Tracks the current state of the application, which determines how to update and display.
	state appState
	// Presents the list of turns from the most recent round of play.
	turnView viewport.Model
	// Tracks the winner from the most recent round of play.
	winner game.Player
}

/*
New creates and returns a model with defaults for a new execution.
*/
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
	helpModel.Width = 36

	viewportModel := viewport.New(74, 50)
	viewportModel.Style = viewportModel.Style.Copy().Padding(1, 2)

	return model{
		bindHelp:  helpModel,
		colorList: colorList,
		game:      game.Game{},
		nameInput: textinput.New(),
		state:     mainState,
		turnView:  viewportModel,
	}
}

/*
Init meets the requirements of the tea.Model interface but is currently a no-op.
*/
func (m model) Init() tea.Cmd {
	return nil
}

/*
Update passes event messages to the current state per requirements of the tea.Model interface.
*/
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Break glass functionality
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		newHeight := sizeMsg.Height - 15
		mainPane = mainPane.Copy().Height(newHeight)
		m.turnView.Height = newHeight
	}

	return m.state.update(msg, m)
}

/*
View renders the current state per requirements of the tea.Model interface.
*/
func (m model) View() string {
	return m.state.view(m)
}

/*
renderHelpContent renders keyMap's view in the help view component of m.
*/
func renderHelpContent(m model, keyMap help.KeyMap) string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		helpTitle,
		m.bindHelp.View(keyMap),
	)
}

/*
renderPlayers takes players and renders it for proper display.
If m has a winner set that winner will be marked as such in the output.
Additionally, if selected is a valid int from the range of players, the corresponding player will be highlighted.
*/
func renderPlayers(m model, selected int) string {
	var (
		players = m.game.Players()
		rows    []string
	)

	for i := 0; i < len(players); i++ {
		var (
			name        string
			playerColor lipgloss.Style
			prefix      string
		)

		switch players[i].Color() {
		case game.Blue:
			playerColor = styleBlue
		case game.Green:
			playerColor = styleGreen
		case game.Red:
			playerColor = styleRed
		case game.Yellow:
			playerColor = styleYellow
		}

		if len(players[i].Name) <= (26 - 3) { // content width after padding, minus the prefix
			name = players[i].Name
		} else {
			name = players[i].Name[:26-5] + "…" // leave two spaces for potential double-width rendering of ellipsis
		}

		if i == selected {
			prefix = " > "
			playerColor = playerColor.Copy().Background(violet)
		} else if m.winner.Name == players[i].Name {
			prefix = "👑 "
		} else {
			prefix = "   "
		}

		rows = append(rows, prefix+playerColor.Render(name))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		playersTitle,
		strings.Join(rows, "\n"))
}

/*
assembleView puts the content pieces together according to the designed layout.
The playersContent, helpContent, and mainContent will each be placed in their appropriate locations in the view.
appStates simply need to determine what they want placed in each pane and pass the content to this function.
*/
func assembleView(playersContent string, helpContent string, mainContent string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainPane.Render(mainContent),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			playersPane.Render(playersContent),
			helpPane.Render(helpContent),
		),
	)
}
