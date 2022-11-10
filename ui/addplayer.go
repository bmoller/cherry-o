package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bmoller/cherry-o/game"
)

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
		cmd  tea.Cmd
		cmds []tea.Cmd
		err  error
	)

	switch keyMsg, ok := msg.(tea.KeyMsg); {
	case ok && keyMsg.Type == tea.KeyEnter:
		// all done; make the call to add a player
		if m.game, err = m.game.AddPlayer(m.nameInput.Value(), m.colorList.SelectedItem().(game.Color)); err != nil {
			m.currentState = errorState
			m.err = err
		} else {
			var players []list.Item
			for _, player := range m.game.Players() {
				players = append(players, player)
			}
			cmds = append(cmds, m.playerList.SetItems(players))
			m.currentState = mainState
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

func viewAddPlayerState(m model) string {
	return ""
}
