package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	rootCommand = &cobra.Command{
		Use:   "blah blah",
		Short: "do the thing",
		RunE:  root,
	}
)

func Execute() error {
	return rootCommand.Execute()
}

func root(cmd *cobra.Command, args []string) error {
	var playerCount int
	var line string

	fmt.Println("Welcome!")
	fmt.Println("How many players?")
	_, err := fmt.Scanln(&line)
	if err != nil {
		log.Println(err)
		log.Fatal("Failed to read from stdin!")
	}
	_, err = fmt.Sscanf(line, "%d", &playerCount)
	for err != nil {
		fmt.Println(err)
		fmt.Printf("%s is not a valid number\n", line)
		fmt.Println("How many players?")
		_, err = fmt.Scanln(&line)
		if err != nil {
			log.Println(err)
			log.Fatal("Failed to read from stdin!")
		}
		_, err = fmt.Sscanf(line, "%d", &playerCount)
	}
	fmt.Printf("You asked for %d players\n", playerCount)

	return nil
}
