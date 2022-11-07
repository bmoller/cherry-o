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
	players map[string]*Player
}

func New() *Game {
	return &Game{
		players: make(map[string]*Player, 4),
	}
}

func (g *Game) Players() map[string]*Player {
	return g.players
}

func (g *Game) AvailableColors() (colors []Color) {
	stateMap := map[Color]bool{
		Blue:   true,
		Green:  true,
		Red:    true,
		Yellow: true,
	}

	for _, player := range g.players {
		stateMap[player.color] = false
	}

	for color, available := range stateMap {
		if available {
			colors = append(colors, color)
		}
	}

	return
}

func (g *Game) AddPlayer(name string, color Color) error {
	if g.players != nil && len(g.players) == 4 {
		return errors.New("max player count reached; unable to add a new player")
	}

	found := false
	for _, availableColor := range g.AvailableColors() {
		if color == availableColor {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("color %s is not available", color)
	}

	g.players[name] = &Player{
		Name:  name,
		color: color,
	}

	return nil
}

func (g *Game) RemovePlayer(player *Player) (err error) {
	if g.players == nil || len(g.players) == 0 {
		err = errors.New("no players to remove")
	} else {
		var found bool
		for name := range g.players {
			if name == player.Name {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("%s is not a current player", player.Name)
		} else {
			delete(g.players, player.Name)
		}
	}

	return
}

func (g *Game) Play() (transcript []Turn, winner *Player, err error) {
	if g.players == nil || len(g.players) == 0 {
		err = errors.New("need at least one player to play")
		return
	}

	var (
		currentTurn Turn
		gameOver    bool
	)
	for !gameOver {
		for _, currentPlayer := range g.players {
			currentTurn, err = takeTurn(currentPlayer)
			if err != nil {
				return
			} else {
				transcript = append(transcript, currentTurn)
			}
			if currentPlayer.cherries == winningScore {
				gameOver = true
				winner = currentPlayer
				break
			}
		}
	}

	return
}

type Turn struct {
	spin   int
	player *Player
}

func (t Turn) Spin() int {
	return t.spin
}

func (t Turn) Player() *Player {
	return t.player
}

func takeTurn(player *Player) (result Turn, err error) {
	upperBound := big.NewInt(int64(len(spinnerValues)))
	spin, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		err = fmt.Errorf("failed to generate a random number: %s", err)
	} else {
		value := spinnerValues[spin.Uint64()]
		player.UpdateCherries(value)
		result = Turn{
			spin:   value,
			player: player,
		}
	}

	return
}
