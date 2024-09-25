package views

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndfsa/cardboard-bank/common/repository"
)

const (
	stateAuth = iota
	stateDashboard
)

type viewState int

var modelStyle = lipgloss.NewStyle()

type MainModel struct {
	width  int
	height int
	state  viewState
	auth   AuthModel
	dash   DashboardModel
}

func NewMainModel(db *sql.DB) MainModel {
	authRepo := repository.NewAuthRepository(db)
	return MainModel{
		state: stateAuth,
		auth:  NewAuthModel(authRepo),
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
	case AuthSuccessMsg:
		m.state = stateDashboard
        m.dash = NewDashboardModel(msg.User)
	}

	switch m.state {
	case stateAuth:
		am, cmd := m.auth.Update(msg)
		m.auth = am.(AuthModel)
		return m, cmd
	case stateDashboard:
		dm, cmd := m.dash.Update(msg)
		m.dash = dm.(DashboardModel)
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
	case stateDashboard:
		return m.dash.View()
	default:
		panic("illegal state")
	}
}
