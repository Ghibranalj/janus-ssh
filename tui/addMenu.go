package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	var b strings.Builder

	// Title
	b.WriteString(Title.Render("Add New Server"))
	b.WriteString("\n\n")

	// Username field
	usernameCursor := " "
	var usernameStyle lipgloss.Style
	if m.activeField == 0 {
		usernameCursor = Cursor.Render(">")
		usernameStyle = FieldActive
	} else {
		usernameStyle = FieldInactive
	}

	usernameLabel := FieldLabel.Render("Username:")
	usernameValue := usernameStyle.Render(m.username)
	b.WriteString(fmt.Sprintf("%s %s %s\n", usernameCursor, usernameLabel, usernameValue))

	// Host field
	hostCursor := " "
	var hostStyle lipgloss.Style
	if m.activeField == 1 {
		hostCursor = Cursor.Render(">")
		hostStyle = FieldActive
	} else {
		hostStyle = FieldInactive
	}

	hostLabel := FieldLabel.Render("Host:")
	hostValue := hostStyle.Render(m.host)
	b.WriteString(fmt.Sprintf("%s %s %s\n", hostCursor, hostLabel, hostValue))

	// Help text
	b.WriteString("\n")
	b.WriteString(HelpText.Render("Tab: Switch fields | Enter: Submit | Esc: Cancel"))

	// Error message
	if m.errorMessage != "" {
		b.WriteString("\n\n")
		b.WriteString(ErrorMsg.Render("Error: " + m.errorMessage))
	}

	return tea.NewView(b.String())
}
