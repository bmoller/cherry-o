package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bmoller/cherry-o/pkg/ui"
)

func main() {
	// playerNames := map[string]game.Color{
	// 	"Brandon":  game.Blue,
	// 	"Madelyn":  game.Green,
	// 	"Geoffrey": game.Red,
	// 	"Austin":   game.Yellow,
	// }
	// var players []*game.Player
	// for name, color := range playerNames {
	// 	newPlayer, err := game.NewPlayer(color)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	newPlayer.Name = name
	// 	players = append(players, newPlayer)
	// }
	// game.NewGame(players)
	// transcript, err := game.Play()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Print(transcript)
	// fmt.Printf("The game lasted %d rounds\n", game.Rounds())
	prog := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if err := prog.Start(); err != nil {
		fmt.Printf("Boo-boo :( - %s\n", err)
		os.Exit(1)
	}
}
