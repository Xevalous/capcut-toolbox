package ui

import (
	"fmt"

	"capcut-toolbox/lock"

	tea "github.com/charmbracelet/bubbletea"
)

type lockModel struct {
	appDir  string
	status  string
	running bool
	done    bool
	msg     string
}

func NewLock(appDir string) tea.Model {
	s := "Unlocked"
	if lock.IsLocked(appDir) {
		s = "Locked"
	}
	return lockModel{appDir: appDir, status: s}
}

func (m lockModel) Init() tea.Cmd {
	return nil
}

func (m lockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.done {
			return NewMenu(m.appDir), nil
		}
		if m.running {
			return m, nil
		}
		switch msg.String() {
		case "1":
			m.running = true
			return m, func() tea.Msg {
				err := lock.On(m.appDir)
				if err != nil {
					return lockResultMsg{err: err}
				}
				return lockResultMsg{}
			}
		case "2":
			m.running = true
			return m, func() tea.Msg {
				err := lock.Off(m.appDir)
				if err != nil {
					return lockResultMsg{err: err}
				}
				return lockResultMsg{}
			}
		case "3":
			return NewMenu(m.appDir), nil
		}
	case lockResultMsg:
		m.running = false
		m.done = true
		if msg.err != nil {
			m.msg = fmt.Sprintf("Error: %v", msg.err)
		} else {
			if lock.IsLocked(m.appDir) {
				m.status = "Locked"
			} else {
				m.status = "Unlocked"
			}
			m.msg = "Done!"
		}
	}
	return m, nil
}

func (m lockModel) View() string {
	s := "\n       LOCK VERSION\n\n"
	s += fmt.Sprintf("  Status: %s\n\n", m.status)

	if m.running {
		s += "  Processing...\n"
	} else if m.done {
		s += fmt.Sprintf("  %s\n\n", m.msg)
		s += "  Press any key to go back...\n"
	} else {
		s += "  [1] Lock\n"
		s += "  [2] Unlock\n"
		s += "  [3] Back\n"
	}

	return s
}

type lockResultMsg struct {
	err error
}
