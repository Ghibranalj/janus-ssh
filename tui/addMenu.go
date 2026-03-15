package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// AddMenu handles adding new servers
type AddMenu struct {
	username     string
	host         string
	activeField  int // 0: username, 1: host
	errorMessage string
}

func NewAddMenu() *AddMenu {
	return &AddMenu{
		activeField: 0,
	}
}

func (m *AddMenu) Init() tea.Cmd {
	return nil
}

func (m *AddMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.username == "" || m.host == "" {
				m.errorMessage = "Both fields are required"
			} else {
				server := Server{User: m.username, Host: m.host}
				return m, func() tea.Msg { return ServerAddedMsg{Server: server} }
			}
		case "tab":
			if m.activeField == 0 {
				m.activeField = 1
			}
		case "shift+tab":
			if m.activeField == 1 {
				m.activeField = 0
			}
		case "esc":
			return m, func() tea.Msg { return SwitchToSelectMsg{} }
		case "backspace":
			if m.activeField == 0 && len(m.username) > 0 {
				m.username = m.username[:len(m.username)-1]
			} else if m.activeField == 1 && len(m.host) > 0 {
				m.host = m.host[:len(m.host)-1]
			}
		default:
			// Handle character input
			keyStr := msg.String()
			if len(keyStr) == 1 && keyStr >= " " && keyStr <= "~" {
				if m.activeField == 0 {
					m.username += keyStr
				} else {
					m.host += keyStr
				}
			}
		}
	}
	return m, nil
}

func (m *AddMenu) View() tea.View {
	var s strings.Builder
	s.WriteString("Add New Server\n\n")

	// Username field
	usernameCursor := " "
	if m.activeField == 0 {
		usernameCursor = ">"
	}
	s.WriteString(fmt.Sprintf("%s Username: %s\n", usernameCursor, m.username))

	// Host field
	hostCursor := " "
	if m.activeField == 1 {
		hostCursor = ">"
	}
	s.WriteString(fmt.Sprintf("%s Host: %s\n", hostCursor, m.host))

	// Help text
	s.WriteString("\n")
	s.WriteString("Tab: Switch fields\n")
	s.WriteString("Enter: Submit\n")
	s.WriteString("Esc: Cancel\n")

	// Error message
	if m.errorMessage != "" {
		s.WriteString(fmt.Sprintf("\nError: %s\n", m.errorMessage))
	}

	return tea.NewView(s.String())
}
