package tui

import (
	"fmt"
	"io"

	tea "charm.land/bubbletea/v2"
	"github.com/gliderlabs/ssh"
)

// Menu transition messages
type SwitchToSelectMsg struct{}
type SwitchToAddMsg struct{}
type ServerSelectedMsg struct{ Server Server }
type ServerAddedMsg struct{ Server Server }

// appModel is the root model that coordinates between menus
type appModel struct {
	currentMenu    tea.Model
	repo           ServerRepository
	selectedServer Server
}

func NewApp(repo ServerRepository) (*appModel, error) {
	return &appModel{
		currentMenu: NewSelectMenu(repo),
		repo:        repo,
	}, nil
}

func (m *appModel) Init() tea.Cmd {
	return nil
}

func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle menu transition messages
	switch msg := msg.(type) {
	case SwitchToSelectMsg:
		m.currentMenu = NewSelectMenu(m.repo)
		return m, nil
	case SwitchToAddMsg:
		m.currentMenu = NewAddMenu()
		return m, nil
	case ServerSelectedMsg:
		m.selectedServer = msg.Server
		return m, tea.Quit
	case ServerAddedMsg:
		// Add to repo
		m.repo.Add(msg.Server)
		m.currentMenu = NewSelectMenu(m.repo)
		return m, nil
	}

	// Delegate update to current menu
	var cmd tea.Cmd
	m.currentMenu, cmd = m.currentMenu.Update(msg)
	return m, cmd
}

func (m *appModel) View() tea.View {
	return m.currentMenu.View()
}

// RunMenu displays the Bubble Tea menu and returns the selected server
func RunMenu(s ssh.Session, repo ServerRepository) (Server, error) {
	// Get PTY information
	ptyReq, winCh, ok := s.Pty()

	// Prepare environment with TERM variable
	// This is critical for Bubble Tea to detect terminal capabilities
	env := s.Environ()
	if ok && ptyReq.Term != "" {
		env = append([]string{fmt.Sprintf("TERM=%s", ptyReq.Term)}, env...)
	}

	// Create program options
	opts := []tea.ProgramOption{
		tea.WithInput(s),
		tea.WithOutput(s),
	}

	// Pass environment so Bubble Tea can detect terminal type
	if len(env) > 0 {
		opts = append(opts, tea.WithEnvironment(env))
	}

	// Pass initial window size since we can't query terminal
	if ok {
		opts = append(opts, tea.WithWindowSize(ptyReq.Window.Width, ptyReq.Window.Height))
	}

	app, err := NewApp(repo)
	if err != nil {
		return Server{}, err
	}

	p := tea.NewProgram(app, opts...)

	// Handle window resize events
	if ok && winCh != nil {
		go func() {
			for win := range winCh {
				p.Send(tea.WindowSizeMsg{Width: win.Width, Height: win.Height})
			}
		}()
	}

	finalModel, err := p.Run()
	if err != nil {
		return Server{}, err
	}

	m := finalModel.(*appModel)
	return m.selectedServer, nil
}

// RestoreTerminal writes ANSI escape sequences to reset terminal state
func RestoreTerminal(w io.Writer) {
	fmt.Fprint(w, "\033[H\033[2J")
}
