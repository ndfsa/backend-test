package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ndfsa/cardboard-bank/tui/views"
)

func main() {
	p := tea.NewProgram(views.NewMainModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
