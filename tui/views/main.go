package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	stateAuth = iota
)

var modelStyle = lipgloss.NewStyle()

type MainModel struct {
	width  int
	height int
	state  int
	auth   AuthModel
}

func NewMainModel() MainModel {
	return MainModel{
		state: stateAuth,
		auth:  NewAuthModel(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.auth.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.auth.width = msg.Width
		m.auth.height = msg.Height
	}
	switch m.state {
	case stateAuth:
		am, cmd := m.auth.Update(msg)
		m.auth = am.(AuthModel)
		return m, cmd
	default:
		panic("illegal state")
	}
}

func (m MainModel) View() string {
	if m.width == 0 {
		return "loading..."
	}

	switch m.state {
	case stateAuth:
		return m.auth.View()
	default:
		panic("illegal state")
	}
}
