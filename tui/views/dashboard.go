package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ndfsa/cardboard-bank/common/model"
)

type DashboardModel struct {
	width       int
	height      int
	currentUser model.User
}

func NewDashboardModel(user model.User) DashboardModel {
	return DashboardModel{
		currentUser: user,
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return nil
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DashboardModel) View() string {
	return "Welcome " + m.currentUser.Fullname
}
