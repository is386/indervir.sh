package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	white = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	gray  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	dim   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type errMsg error
type loadingDoneMsg struct{}

type navItem struct {
	title string
	color string
}

type model struct {
	style       lipgloss.Style
	spinner     spinner.Model
	navItems    []navItem
	navSelected int
	width       int
	height      int
	loading     bool
	quitting    bool
	err         error
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{" ", "█"},
		FPS:    time.Second / 2,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))

	return model{
		spinner: s,
		loading: false,
		navItems: []navItem{
			{title: "about", color: "34"},
			{title: "coding", color: "205"},
			{title: "gaming", color: "220"},
			{title: "reading", color: "141"},
			{title: "running", color: "208"},
		},
	}
}

// Runs once per start up
func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return loadingDoneMsg{}
	}))
}

// Runs on every event (keypress, window resize, etc)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.navSelected < len(m.navItems)-1 {
				m.navSelected++
			}
			return m, nil
		case "k", "up":
			if m.navSelected > 0 {
				m.navSelected--
			}
			return m, nil
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
	case loadingDoneMsg:
		m.loading = false
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// Renders to screen
func (m model) View() tea.View {
	if m.err != nil {
		v := tea.NewView(m.err.Error())
		v.AltScreen = true
		return v
	}

	var str string
	if m.loading {
		str = m.showLoadingScreen()
	} else {
		str = m.showMainScreen()
	}

	if m.quitting {
		v := tea.NewView(str + "\n")
		v.AltScreen = true
		return v
	}

	v := tea.NewView(str)
	v.AltScreen = true
	return v
}

func (m model) showLoadingScreen() string {
	return lipgloss.NewStyle().
		Width(m.width).
		AlignHorizontal(lipgloss.Center).
		Height(m.height).
		AlignVertical(lipgloss.Center).
		Render(fmt.Sprintf(
			"%s %s",
			white.Bold(true).Render("indervir.dev"),
			m.spinner.View(),
		))
}

func (m model) showMainScreen() string {
	if m.width < 30 || m.height < 10 {
		return "Terminal too small. Please resize."
	}

	maxWidth := 80
	maxHeight := 24
	w := m.width
	h := m.height
	if w > maxWidth {
		w = maxWidth
	}
	if h > maxHeight {
		h = maxHeight
	}

	innerWidth := w - 2
	innerHeight := h - 2
	navWidth := innerWidth / 4
	contentWidth := innerWidth - navWidth - 1

	navbar := lipgloss.NewStyle().
		Width(navWidth).
		Height(innerHeight).
		AlignHorizontal(lipgloss.Left).
		PaddingLeft(1).
		Render(m.renderNav())

	dividerLine := strings.TrimRight(strings.Repeat("│\n", innerHeight), "\n")
	divider := dim.
		Height(innerHeight).
		Render(dividerLine)

	content := m.renderContent(contentWidth, innerHeight)

	inner := lipgloss.JoinHorizontal(lipgloss.Center, navbar, divider, content)

	box := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(dim.GetForeground()).
		Width(w - 2).
		Height(h - 3).
		Render(inner)

	help := dim.
		Width(w - 2).
		AlignHorizontal(lipgloss.Center).
		Render("↑/k up • ↓/j down • q quit")

	page := lipgloss.JoinVertical(lipgloss.Center, box, help)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(page)
}

func (m model) renderNav() string {
	var items []string
	navTitle := white.Bold(true).Render("</> indervir.dev")
	items = append(items, navTitle)
	items = append(items, "")
	for i, nav := range m.navItems {
		if i == m.navSelected {
			items = append(items, lipgloss.NewStyle().
				Foreground(lipgloss.Color(nav.color)).
				Render("❯ "+nav.title))
		} else {
			items = append(items, dim.Render("  "+nav.title))
		}
	}
	return strings.Join(items, "\n")
}

func (m model) renderContent(contentWidth int, innerHeight int) string {
	navItem := m.navItems[m.navSelected]
	title := m.renderContentTitle(contentWidth, navItem)

	var body string
	switch navItem.title {
	case "about":
		body = m.renderAbout(contentWidth)
	case "coding":
		body = m.renderCoding()
	case "gaming":
		body = m.renderGaming()
	case "reading":
		body = m.renderReading()
	case "running":
		body = m.renderRunning()
	}

	wrappedBody := lipgloss.NewStyle().
		Width(contentWidth - 4).
		Render(body)
	content := title + "\n" + wrappedBody
	return lipgloss.NewStyle().
		Width(contentWidth).
		MaxWidth(contentWidth).
		Height(innerHeight).
		PaddingLeft(1).
		Render(content)
}

func (m model) renderContentTitle(contentWidth int, item navItem) string {
	title := lipgloss.NewStyle().
		Width(contentWidth - 4).
		Foreground(lipgloss.Color(item.color)).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Render(item.title)
	line := dim.Render(strings.Repeat("─", contentWidth-4))
	return title + "\n" + line
}

func (m model) renderAbout(contentWidth int) string {
	divider := dim.Render(strings.Repeat("─", contentWidth-4))

	bio := white.Render("\nhey there - my name is indervir singh. i am a software developer, reader, gamer, and very slow runner.\n")
	bio += gray.Render("i was randomly inspired to make this page after discovering terminal.shop. its written in golang, using the charm suite of terminal ui libraries.\n")

	info := strings.Join([]string{
		m.renderInfoRow("location", "new jersey, usa"),
		m.renderInfoRow("contact", "singh.indervir89@gmail.com"),
		m.renderInfoRow("github", lipgloss.NewStyle().Hyperlink("https://github.com/is386").Render("is386")),
		m.renderInfoRow("watch this", lipgloss.NewStyle().Hyperlink("https://youtu.be/gKQOXYB2cd8?si=lmvBPGsDfdDW5LZ-").Render("youtube video")),
	}, "\n\n")

	return bio + divider + "\n\n" + info
}

func (m model) renderCoding() string {
	return "coming soon..."
}

func (m model) renderGaming() string {
	return "coming soon..."
}

func (m model) renderReading() string {
	return "coming soon..."
}

func (m model) renderRunning() string {
	return "coming soon..."
}

func (m model) renderInfoRow(label string, value string) string {
	return gray.Render(label+" ") + white.Render(value)
}
