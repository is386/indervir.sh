package tui

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	style    lipgloss.Style
	spinner  spinner.Model
	width    int
	height   int
	quitting bool
	err      error
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{" ", "█"},
		FPS:    time.Second / 2,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
	return model{spinner: s}
}

// Runs once per start up
func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Runs on every event (keypress, window resize, etc)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

type errMsg error

// Renders to screen
func (m model) View() tea.View {
	if m.err != nil {
		v := tea.NewView(m.err.Error())
		v.AltScreen = true
		return v
	}
	str := lipgloss.NewStyle().
		Width(m.width).
		AlignHorizontal(lipgloss.Center).
		Height(m.height).
		AlignVertical(lipgloss.Center).
		Render(fmt.Sprintf(
			"%s %s",
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255")).Render("indervir.dev"),
			m.spinner.View(),
		))
	if m.quitting {
		v := tea.NewView(str + "\n")
		v.AltScreen = true
		return v
	}
	v := tea.NewView(str)
	v.AltScreen = true
	return v
}
