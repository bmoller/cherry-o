package game

import (
	"fmt"
)

type Color int

const (
	InvalidColor Color = iota
	Blue
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

func (c Color) FilterValue() string {
	return c.String()
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
	case p.cherries > WinningScore:
		p.cherries = WinningScore
	}

	return p
}

func (p Player) Color() Color {
	return p.color
}

func (p Player) String() string {
	return p.Name
}

func (p Player) FilterValue() string {
	return p.String()
}
