package game

import (
	"fmt"
)

type Color int

const (
	Blue Color = iota
	Green
	Red
	Yellow
)

func (c Color) String() string {
	switch c {
	case Blue:
		return "blue"
	case Green:
		return "green"
	case Red:
		return "red"
	case Yellow:
		return "yellow"
	default:
		return fmt.Sprintf("%d", int(c))
	}
}

type Player struct {
	Name string

	cherries int
	color    Color
}

func (player *Player) Cherries() int {
	return player.cherries
}

func (player *Player) UpdateCherries(count int) {
	player.cherries += count
	if player.cherries > winningScore {
		player.cherries = winningScore
	} else if player.cherries < 0 {
		player.cherries = 0
	}
}

func (player *Player) Color() Color {
	return player.color
}

func (player *Player) String() string {
	return player.Name
}
