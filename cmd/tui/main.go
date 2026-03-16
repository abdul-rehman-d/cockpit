package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/abdul-rehman-d/cockpit/internal/view"
)

func main() {
	if _, err := tea.NewProgram(view.NewModel()).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
