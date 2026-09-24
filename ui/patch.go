package ui

import (
	"fmt"
	"path/filepath"

	"capcut-toolbox/patch"
	"capcut-toolbox/util"

	tea "github.com/charmbracelet/bubbletea"
)

type patchModel struct {
	appsDir   string
	versions  []string
	cursor    int
	selected  string
	exeDir    string
	status    string
	running   bool
	done      bool
	msg       string
	selecting bool
}

func NewPatch(appsDir string) tea.Model {
	versions := util.ListVersions(appsDir)
	if len(versions) == 0 {
		return patchModel{appsDir: appsDir, msg: "No versions found under Apps/"}
	}
	if len(versions) == 1 {
		exeDir := filepath.Join(appsDir, versions[0])
		return patchModel{appsDir: appsDir, versions: versions, selected: versions[0], exeDir: exeDir, status: patchStatus(exeDir)}
	}
	return patchModel{appsDir: appsDir, versions: versions, selecting: true}
}

func patchStatus(exeDir string) string {
	if patch.IsPatched(exeDir) {
		return "On"
	}
	return "Off"
}

func (m patchModel) Init() tea.Cmd {
	return nil
}

func (m patchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.done {
			return NewMenu(m.appsDir), nil
		}
		if m.running {
			return m, nil
		}
		if m.selecting {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.versions)-1 {
					m.cursor++
				}
			case "enter":
				m.selected = m.versions[m.cursor]
				m.exeDir = filepath.Join(m.appsDir, m.selected)
				m.status = patchStatus(m.exeDir)
				m.selecting = false
			case "3", "esc":
				return NewMenu(m.appsDir), nil
			}
			return m, nil
		}
		switch msg.String() {
		case "1":
			m.running = true
			return m, func() tea.Msg {
				err := patch.On(m.exeDir)
				if err != nil {
					return patchResultMsg{err: err}
				}
				return patchResultMsg{}
			}
		case "2":
			m.running = true
			return m, func() tea.Msg {
				err := patch.Off(m.exeDir)
				if err != nil {
					return patchResultMsg{err: err}
				}
				return patchResultMsg{}
			}
		case "3":
			return NewMenu(m.appsDir), nil
		}
	case patchResultMsg:
		m.running = false
		m.done = true
		if msg.err != nil {
			m.msg = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.status = patchStatus(m.exeDir)
			m.msg = "Done!"
		}
	}
	return m, nil
}

func (m patchModel) View() string {
	if m.msg == "No versions found under Apps/" {
		return "\n          PATCH\n\n  " + m.msg + "\n\n  Press any key to go back...\n"
	}

	if m.selecting {
		s := "\n          PATCH\n\n  Select version:\n\n"
		for i, v := range m.versions {
			cursor := "  "
			if i == m.cursor {
				cursor = ">>"
			}
			s += fmt.Sprintf("  %s %s\n", cursor, v)
		}
		s += "\n  [enter] Select  [3] Back\n"
		return s
	}

	s := fmt.Sprintf("\n          PATCH\n\n  Version: %s\n  Status: %s\n\n", m.selected, m.status)

	if m.running {
		s += "  Processing...\n"
	} else if m.done {
		s += fmt.Sprintf("  %s\n\n", m.msg)
		s += "  Press any key to go back...\n"
	} else {
		s += "  [1] Patch ON\n"
		s += "  [2] Patch OFF\n"
		s += "  [3] Back\n"
	}

	return s
}

type patchResultMsg struct {
	err error
}
