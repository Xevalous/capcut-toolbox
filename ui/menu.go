package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type menuModel struct {
	appsDir string
}

func NewMenu(appsDir string) tea.Model {
	return menuModel{appsDir: appsDir}
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			return NewPatch(m.appsDir), nil
		case "2":
			return NewLock(m.appsDir), nil
		case "3":
			return NewDownload(m.appsDir), nil
		case "4", "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	s := "\n          CAPCUT TOOLBOX\n\n"
	s += "  [1] Patch\n"
	s += "  [2] Lock Version\n"
	s += "  [3] Download Supported Version\n"
	s += "  [4] Exit\n\n"
	s += "  Created by Xevalous, Visit the repo:\n"
	s += "  https://github.com/Xevalous/capcut-toolbox\n"
	return s
}
