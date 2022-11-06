package game

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
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
	gameOver      bool
	players       []*Player
	rounds        int
	spinnerValues = [7]int{1, 2, 3, 4, -2, -2, -10}
	winner        *Player
)

func GameOver() bool {
	return gameOver
}

func NewGame(group []*Player) {
	gameOver = false
	players = group
	rounds = 0
}

func Rounds() int {
	return rounds
}

func Winner() (*Player, error) {
	if !gameOver {
		return nil, errors.New("game has not ended")
	}

	return winner, nil
}

func Play() (transcript string, err error) {
	if gameOver || players == nil || rounds > 0 {
		err = errors.New("start a new game before beginning to play")
		return
	}

	var (
		textColor string
	)

	for !gameOver {
		for _, currentPlayer := range players {
			summary := takeTurn(currentPlayer)
			switch currentPlayer.color {
			case Blue:
				textColor = terminalBlue
			case Green:
				textColor = terminalGreen
			case Red:
				textColor = terminalRed
			case Yellow:
				textColor = terminalYellow
			default:
				err = fmt.Errorf("player %s has unrecognized color", currentPlayer.Name)
				return
			}
			transcript += fmt.Sprintf("%s%s%s\n", textColor, summary, terminalReset)
			if currentPlayer.cherries == winningScore {
				winner = currentPlayer
				gameOver = true
				transcript += fmt.Sprintf("%s%s won!%s\n", textColor, currentPlayer.Name, terminalReset)
				break
			}
		}
	}

	return
}

func takeTurn(player *Player) (summary string) {
	upperBound := big.NewInt(int64(len(spinnerValues)))
	currentSpin, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		log.Fatal(err)
	}
	value := spinnerValues[currentSpin.Uint64()]

	if value < 0 {
		value *= -1
		player.RemoveCherries(value)
		summary = fmt.Sprintf("Oh noes! %s lost %d cherries!!!", player.Name, currentSpin)
	} else {
		player.AddCherries(value)
		switch value {
		case 1:
			summary = fmt.Sprintf("Yay! %s got 1 more cherry!!!", player.Name)
		default:
			summary = fmt.Sprintf("Yay! %s got %d more cherries!!!", player.Name, currentSpin)
		}
	}

	return
}
