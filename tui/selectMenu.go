package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// SelectMenu handles server selection
type SelectMenu struct {
	cursor int
	repo   ServerRepository
}

func NewSelectMenu(repo ServerRepository) *SelectMenu {
	return &SelectMenu{
		repo: repo,
	}
}

func (m *SelectMenu) Init() tea.Cmd {
	return nil
}

func (m *SelectMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "a":
			return m, func() tea.Msg { return SwitchToAddMsg{} }
		case "up":
			servers, _ := m.repo.List()
			if m.cursor > 0 {
				m.cursor--
			} else if len(servers) > 0 {
				m.cursor = len(servers) - 1
			}
		case "down":
			servers, _ := m.repo.List()
			if m.cursor < len(servers)-1 {
				m.cursor++
			} else if len(servers) > 0 {
				m.cursor = 0
			}
		case "enter":
			servers, _ := m.repo.List()
			if m.cursor >= 0 && m.cursor < len(servers) {
				server := servers[m.cursor]
				return m, func() tea.Msg { return ServerSelectedMsg{Server: Server{User: server.User, Host: server.Host}} }
			}
		}
	}
	return m, nil
}

func (m *SelectMenu) View() tea.View {
	servers, _ := m.repo.List()
	v := "Select server:\n\n"
	for i, server := range servers {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		v += fmt.Sprintf("%s %s\n", cursor, server.String())
	}
	v += "\nq to quit | a: Add new server\n"
	return tea.NewView(v)
}
