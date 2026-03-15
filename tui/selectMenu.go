package tui

import (
	"fmt"
	"strings"

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

	var b strings.Builder

	// Title
	b.WriteString(Title.Render("Select server"))
	b.WriteString("\n\n")

	// Server list
	for i, server := range servers {
		cursor := " "
		if m.cursor == i {
			cursor = Cursor.Render(">")
		}

		var itemText string
		if m.cursor == i {
			itemText = SelectedItem.Render(server.String())
		} else {
			itemText = NormalItem.Render(server.String())
		}

		b.WriteString(fmt.Sprintf("%s %s\n", cursor, itemText))
	}

	// Help text
	b.WriteString("\n")
	b.WriteString(HelpText.Render("q: Quit | a: Add new server"))

	return tea.NewView(b.String())
}
