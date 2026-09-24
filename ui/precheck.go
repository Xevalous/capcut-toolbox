package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"capcut-toolbox/util"

	tea "github.com/charmbracelet/bubbletea"
)

type checkResult struct {
	name   string
	status string // "OK", "WARN", "FAIL"
}

type precheckModel struct {
	checks  []checkResult
	failed  bool
	appsDir string
}

type tickMsg time.Time

func runPrechecks() []checkResult {
	var checks []checkResult

	// CapCut installed
	localAppData := os.Getenv("LOCALAPPDATA")
	capcutPath := filepath.Join(localAppData, "CapCut")
	if _, err := os.Stat(capcutPath); err == nil {
		checks = append(checks, checkResult{"CapCut installed", "OK"})
	} else {
		checks = append(checks, checkResult{"CapCut installed", "FAIL"})
	}

    // Version check — supported: 1.0.0 ~ 5.4.0
	appsDir, _ := util.ResolveAppsDir()
	if appsDir != "" {
		versions := util.ListVersions(appsDir)
		sort.Strings(versions)
		if len(versions) == 0 {
			checks = append(checks, checkResult{"Version: none installed", "WARN"})
		} else {
			for _, v := range versions {
				status := "OK"
				if !isVersionSupported(v) {
					status = "WARN"
				}
				checks = append(checks, checkResult{"Version: " + v, status})
			}
		}
	}

	// CapCut running
	if util.IsCapCutRunning() {
		checks = append(checks, checkResult{"CapCut is running", "WARN"})
	} else {
		checks = append(checks, checkResult{"CapCut is not running", "OK"})
	}

	return checks
}

func isVersionSupported(ver string) bool {
	parts := strings.Split(ver, ".")
	if len(parts) < 2 {
		return false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	// ponytail: hardcoded range 1.0 ~ 5.4, matches download list
	if major < 1 || major > 5 {
		return false
	}
	if major == 5 && minor > 4 {
		return false
	}
	return true
}

func NewPrecheck() tea.Model {
	checks := runPrechecks()
	failed := false
	for _, c := range checks {
		if c.status == "FAIL" {
			failed = true
			break
		}
	}

	appsDir, _ := util.ResolveAppsDir()

	return precheckModel{checks: checks, failed: failed, appsDir: appsDir}
}

func (m precheckModel) Init() tea.Cmd {
	if m.failed {
		return nil
	}
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m precheckModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		if m.failed {
			return m, tea.Quit
		}
	case tickMsg:
		if !m.failed {
			return NewMenu(m.appsDir), nil
		}
	}
	return m, nil
}

func (m precheckModel) View() string {
	s := "\n  CAPCUT TOOLBOX\n\n  Checking requirements...\n\n"

	for _, c := range m.checks {
		s += fmt.Sprintf("  [%s]  %s\n", c.status, c.name)
	}

	s += "\n"

	if m.failed {
		s += "  Press any key to exit...\n"
	}

	return s
}


