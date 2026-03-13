package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Style lipgloss.Style
}

// Runs once per start up
func (m Model) Init() tea.Cmd {
	return nil
}

// Runs on every event (keypress, window resize, etc)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// Renders to screen
func (m Model) View() string {
	return m.Style.Render("Hello from indervir.sh! Press 'q' to quit.")
}
