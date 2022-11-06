package game

import (
	"errors"
)

type Color int

const (
	Blue Color = iota
	Green
	Red
	Yellow
)

type Player struct {
	Name  string
	Score int

	cherries int
	color    Color
}

func NewPlayer(playerColor Color) (player *Player, err error) {
	player = new(Player)

	if playerColor != Blue && playerColor != Green && playerColor != Red && playerColor != Yellow {
		err = errors.New("invalid player color")
	} else {
		player.color = playerColor
	}

	return
}

func (player *Player) Cherries() int {
	return player.cherries
}

func (player *Player) AddCherries(count int) {
	player.cherries += count
	if player.cherries > winningScore {
		player.cherries = winningScore
	}
}

func (player *Player) RemoveCherries(count int) {
	player.cherries -= count
	if player.cherries < 0 {
		player.cherries = 0
	}
}

func (player *Player) Color() Color {
	return player.color
}

func (player *Player) String() string {
	return player.Name
}
