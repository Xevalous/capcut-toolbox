package main

import (
	"fmt"
	"os"

	"capcut-toolbox/ui"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("capcut-toolbox %s\n", version)
		os.Exit(0)
	}

	p := tea.NewProgram(ui.NewPrecheck())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
