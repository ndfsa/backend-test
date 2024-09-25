package views

import (
	"context"
	"log"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/ndfsa/cardboard-bank/common/repository"
)

type AuthModel struct {
	focusIdx      int
	width         int
	height        int
	userTextInput textinput.Model
	passTextInput textinput.Model
	cursorMode    cursor.Mode
	repo          repository.AuthRepository
}

type AuthSuccessMsg struct {
	User model.User
}

var focusedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("205"))
var defaultStyle = lipgloss.NewStyle()

func NewAuthModel(repo repository.AuthRepository) AuthModel {
	textInputUsername := textinput.New()
	textInputUsername.Placeholder = "Username"
	textInputUsername.CharLimit = 40

	textInputUsername.Focus()
	textInputUsername.TextStyle = focusedStyle
	textInputUsername.PromptStyle = focusedStyle

	textInputPassword := textinput.New()
	textInputPassword.Placeholder = "Password"
	textInputPassword.EchoMode = textinput.EchoPassword
	textInputPassword.EchoCharacter = '*'
	textInputPassword.CharLimit = 40

	m := AuthModel{
		userTextInput: textInputUsername,
		passTextInput: textInputPassword,
		repo:          repo,
	}

	return m
}

func (m AuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			m.focusIdx++

			if m.focusIdx > 1 {
				username := m.userTextInput.Value()
				password := m.passTextInput.Value()
				return m, m.login(username, password)
			}

			switch m.focusIdx {
			case 0:
				m.passTextInput.Blur()
				m.passTextInput.TextStyle = defaultStyle
				m.passTextInput.PromptStyle = defaultStyle
				m.userTextInput.TextStyle = focusedStyle
				m.userTextInput.PromptStyle = focusedStyle
				return m, m.userTextInput.Focus()
			case 1:
				m.userTextInput.Blur()
				m.userTextInput.TextStyle = defaultStyle
				m.passTextInput.TextStyle = focusedStyle
				m.userTextInput.PromptStyle = defaultStyle
				m.passTextInput.PromptStyle = focusedStyle
				return m, m.passTextInput.Focus()
			default:
				m.userTextInput.TextStyle = defaultStyle
				m.passTextInput.TextStyle = defaultStyle
				m.userTextInput.PromptStyle = defaultStyle
				m.passTextInput.PromptStyle = defaultStyle
				m.userTextInput.Blur()
				m.passTextInput.Blur()
			}
		}

	case error:
		log.Println(msg)
		return m, tea.Quit
	case model.User:
		return m, func() tea.Msg {
			return AuthSuccessMsg{User: msg}
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m AuthModel) login(username, password string) func() tea.Msg {
	return func() tea.Msg {
		ctx := context.Background()
		user, err := m.repo.Authenticate(ctx, username, password)
		if err != nil {
			return err
		}

		return user
	}
}

func (m *AuthModel) updateInputs(msg tea.Msg) tea.Cmd {
	var uCmd, pCmd tea.Cmd
	m.userTextInput, uCmd = m.userTextInput.Update(msg)
	m.passTextInput, pCmd = m.passTextInput.Update(msg)

	return tea.Batch(uCmd, pCmd)
}

func (m AuthModel) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		0.2,
		lipgloss.Center,
		modelStyle.Render(lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().
				Width(42).
				Render(m.userTextInput.View()),
			lipgloss.NewStyle().
				Width(42).
				Render(m.passTextInput.View()))))
}
