package game

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

const (
	terminalBlue   = "\u001b[34m"
	terminalGreen  = "\u001b[32m"
	terminalRed    = "\u001b[31m"
	terminalReset  = "\u001b[0m"
	terminalYellow = "\u001b[33m"
	winningScore   = 10
)

var (
	spinnerValues = [7]int{1, 2, 3, 4, -2, -2, -10}
)

type Game struct {
	playerCount int
	players     [4]Player
}

func (g Game) Players() []Player {
	return g.players[:g.playerCount]
}

func (g Game) AvailableColors() (availability map[Color]bool) {
	availability = map[Color]bool{
		Blue:   true,
		Green:  true,
		Red:    true,
		Yellow: true,
	}

	for i := 0; i < g.playerCount; i++ {
		availability[g.players[i].color] = false
	}

	return
}

func (g Game) AddPlayer(name string, color Color) (Game, error) {
	var err error

	switch {
	case g.playerCount == 4:
		err = errors.New("max player count reached; unable to add a new player")
	case !g.AvailableColors()[color]:
		err = fmt.Errorf("the %s color is not available", color)
	default:
		g.players[g.playerCount+1] = Player{
			Name:  name,
			color: color,
		}
		g.playerCount++
	}

	return g, err
}

func (g Game) RemovePlayer(name string) (Game, error) {
	var err error

	switch {
	case g.playerCount == 0:
		err = errors.New("no players to remove")
	case name == "":
		err = errors.New("must provide a valid name")
	default:
		var (
			found   bool
			players [4]Player = [4]Player{}
		)

		for i := 0; i < g.playerCount; i++ {
			if g.players[i].Name == name {
				found = true
				g.playerCount--
			} else {
				players[i] = g.players[i]
			}
		}
		g.players = players

		if !found {
			err = fmt.Errorf("%s is not a current player", name)
		}
	}

	return g, err
}

func (g Game) Play() (turns []Turn, winner Player, err error) {
	switch g.playerCount {
	case 0:
		err = errors.New("need at least one player to play")
	default:
		var (
			turn     Turn
			gameOver bool
			i        int
		)

		for !gameOver {
			for i = 0; i < g.playerCount; i++ {
				turn, g.players[i], err = takeTurn(g.players[i])
				if err != nil {
					return
				} else {
					turns = append(turns, turn)
				}
				if g.players[i].cherries == winningScore {
					gameOver = true
					winner = g.players[i]
					break
				}
			}
		}
	}

	return
}

type Turn struct {
	Spin   int
	Player Player
}

func takeTurn(player Player) (Turn, Player, error) {
	var (
		err  error
		turn Turn
	)

	upperBound := big.NewInt(int64(len(spinnerValues)))
	spin, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		err = fmt.Errorf("failed to generate a random number: %s", err)
	} else {
		value := spinnerValues[spin.Uint64()]
		player = player.updateCherries(value)
		turn = Turn{
			Spin:   value,
			Player: player,
		}
	}

	return turn, player, err
}
