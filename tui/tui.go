package tui

import (
	"fmt"
	"io"

	tea "charm.land/bubbletea/v2"
)

// Menu transition messages
type SwitchToSelectMsg struct{}
type SwitchToAddMsg struct{}
type SwitchToEditMsg struct{ Server Server }
type ServerSelectedMsg struct{ Server Server }
type ServerAddedMsg struct{ Server Server }
type ServerEditedMsg struct {
	OldUser string
	OldHost string
	Server  Server
}
type ServerDeletedMsg struct {
	User string
	Host string
}

type exitTUI struct{}

type sshErrorMsg struct{ err error }

// appModel is the root model that coordinates between menus
type AppModel struct {
	currentMenu    tea.Model
	repo           ServerRepository

	// output
	Server string
	Exit   bool
}

func NewApp(repo ServerRepository) *AppModel {
	return &AppModel{
		currentMenu: NewSelectMenu(repo),
		repo:        repo,
	}
}

func (m *AppModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle menu transition messages
	switch msg := msg.(type) {
	case SwitchToSelectMsg:
		m.currentMenu = NewSelectMenu(m.repo)
		return m, nil
	case SwitchToAddMsg:
		m.currentMenu = NewServerForm()
		return m, nil
	case SwitchToEditMsg:
		m.currentMenu = NewServerFormWithValues(msg.Server.User, msg.Server.Host)
		return m, nil
	case ServerSelectedMsg:
		m.Server = msg.Server.String()
		return m, tea.Quit
	case ServerAddedMsg:
		m.repo.Add(msg.Server)
		m.currentMenu = NewSelectMenu(m.repo)
		return m, nil
	case ServerEditedMsg:
		// Update in repo
		m.repo.Update(msg.OldUser, msg.OldHost, msg.Server)
		m.currentMenu = NewSelectMenu(m.repo)
		return m, nil
	case ServerDeletedMsg:
		// Delete from repo
		m.repo.Delete(msg.User, msg.Host)
		m.currentMenu = NewSelectMenu(m.repo)
		return m, nil
	case sshErrorMsg:
		// Show error and return to menu
		m.currentMenu = NewSelectMenu(m.repo)
		return m, nil
	case exitTUI:
		m.Exit = true
		return m, tea.Quit
	}

	// Delegate update to current menu
	var cmd tea.Cmd
	m.currentMenu, cmd = m.currentMenu.Update(msg)
	return m, cmd
}

func (m *AppModel) View() tea.View {
	return m.currentMenu.View()
}

// RestoreTerminal writes ANSI escape sequences to reset terminal state
func RestoreTerminal(w io.Writer) {
	fmt.Fprint(w, "\033[H\033[2J")
}
