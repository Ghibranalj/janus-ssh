package tui

import (
	"fmt"
	"io"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/ssh"
	"github.com/creack/pty"
)

// Menu transition messages
type SwitchToSelectMsg struct{}
type SwitchToAddMsg struct{}
type ServerSelectedMsg struct{ Server Server }
type ServerAddedMsg struct{ Server Server }

type sshErrorMsg struct{ err error }

// appModel is the root model that coordinates between menus
type appModel struct {
	currentMenu    tea.Model
	repo           ServerRepository
	selectedServer Server
	session        ssh.Session
}

func NewApp(repo ServerRepository, session ssh.Session) *appModel {
	return &appModel{
		currentMenu: NewSelectMenu(repo),
		repo:        repo,
		session:     session,
	}
}

func (m *appModel) Init() tea.Cmd {
	m.ClearTerminal()
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
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
		// Launch SSH session and return to menu when done
		return m, m.launchSSHSession(msg.Server)
	case ServerAddedMsg:
		// Add to repo
		m.repo.Add(msg.Server)
		m.currentMenu = NewSelectMenu(m.repo)
		return m, nil
	case sshErrorMsg:
		// Show error and return to menu
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

func (m *appModel) launchSSHSession(server Server) tea.Cmd {
	return func() tea.Msg {
		m.ClearTerminal()

		ptyReq, winCh, _ := m.session.Pty()

		cmd := exec.Command("ssh", server.String())
		cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))

		ptmx, err := pty.Start(cmd)
		if err != nil {
			return sshErrorMsg{err: err}
		}
		defer ptmx.Close()

		if winCh != nil {
			go func() {
				for win := range winCh {
					pty.Setsize(ptmx, &pty.Winsize{
						Rows: uint16(win.Height),
						Cols: uint16(win.Width),
					})
				}
			}()
		}

		go io.Copy(ptmx, m.session)
		io.Copy(m.session, ptmx)

		cmd.Wait()

		m.ClearTerminal()
		// Return to menu after SSH exits
		return SwitchToSelectMsg{}
	}
}

func (m *appModel) ClearTerminal() {
	RestoreTerminal(m.session)
}

// RestoreTerminal writes ANSI escape sequences to reset terminal state
func RestoreTerminal(w io.Writer) {
	fmt.Fprint(w, "\033[H\033[2J")
}
