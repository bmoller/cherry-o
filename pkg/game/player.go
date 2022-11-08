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

func (p Player) updateCherries(amount int) Player {
	p.cherries += amount

	switch {
	case p.cherries < 0:
		p.cherries = 0
	case p.cherries > 10:
		p.cherries = 10
	}

	return p
}

func (player Player) Color() Color {
	return player.color
}

func (player Player) String() string {
	return player.Name
}
