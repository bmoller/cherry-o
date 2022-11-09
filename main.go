package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bmoller/cherry-o/ui"
)

func main() {
	prog := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if err := prog.Start(); err != nil {
		fmt.Printf("Boo-boo :( - %s\n", err)
		os.Exit(1)
	}
}
