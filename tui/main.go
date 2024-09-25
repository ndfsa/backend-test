package main

import (
	"database/sql"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/cardboard-bank/tui/views"
)

func main() {
	db, err := sql.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}

	p := tea.NewProgram(views.NewMainModel(db), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
