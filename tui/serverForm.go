package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ServerForm handles adding and editing servers
type ServerForm struct {
	username     string
	host         string
	activeField  int // 0: username, 1: host
	errorMessage string
	// Editing state
	isEditing bool   // true if editing an existing server
	oldUser   string // original username (for editing)
	oldHost   string // original host (for editing)
}

func NewServerForm() *ServerForm {
	return &ServerForm{
		activeField: 0,
		isEditing:   false,
	}
}

func NewServerFormWithValues(username, host string) *ServerForm {
	return &ServerForm{
		username:    username,
		host:        host,
		activeField: 0,
		isEditing:   true,
		oldUser:     username,
		oldHost:     host,
	}
}

func (m *ServerForm) Init() tea.Cmd {
	return nil
}

func (m *ServerForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	key := keyMsg.String()

	switch key {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		if !m.isValid() {
			return m, nil
		}
		return m, m.submitMessage()
	case "tab":
		m.activeField++
		m.activeField %= 2
	case "esc":
		return m, func() tea.Msg { return SwitchToSelectMsg{} }
	case "backspace":
		m.handleBackspace()
	default:
		m.handleCharacterInput(key)
	}

	return m, nil
}

func (m *ServerForm) handleBackspace() {
	if m.activeField == 0 && len(m.username) > 0 {
		m.username = m.username[:len(m.username)-1]
	} else if m.activeField == 1 && len(m.host) > 0 {
		m.host = m.host[:len(m.host)-1]
	}
}

func (m *ServerForm) handleCharacterInput(key string) {
	if len(key) != 1 || key < " " || key > "~" {
		return
	}
	if m.activeField == 0 {
		m.username += key
	} else {
		m.host += key
	}
}

func (m *ServerForm) View() tea.View {
	var b strings.Builder

	// Title
	title := "Add New Server"
	if m.isEditing {
		title = "Edit Server"
	}
	b.WriteString(Title.Render(title))
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

func (m *ServerForm) isValid() bool {
	m.errorMessage = ""
	if m.username == "" || m.host == "" {
		m.errorMessage = "Both fields are required"
		return false
	}
	return true
}

func (m *ServerForm) submitMessage() tea.Cmd {
	server := Server{User: m.username, Host: m.host}
	if m.isEditing {
		return func() tea.Msg { return ServerEditedMsg{OldUser: m.oldUser, OldHost: m.oldHost, Server: server} }
	}
	return func() tea.Msg { return ServerAddedMsg{Server: server} }
}
