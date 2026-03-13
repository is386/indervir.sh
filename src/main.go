package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/is386/indervir.sh/src/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
)

const (
	host = "0.0.0.0"
	port = 2235
)

func main() {
	// Generates a new SSH server with the given address, host key, and middleware
	// First, activeterm rejects connections with out a PTY (terminal)
	// Second, bubbletea launches a Bubble Tea TUI app for the SSH session
	wishServer, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		os.Exit(1)
	}

	// Starts the server inside of a go routine and listens for SIGINT/SIGTERM (ctrl+c)
	serverDone := make(chan os.Signal, 1)
	signal.Notify(serverDone, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err := wishServer.ListenAndServe(); err != nil {
			log.Error("Could not start server", "error", err)
			serverDone <- nil
		}
	}()

	// Graceful Exit
	<-serverDone
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := wishServer.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}

// Creates a Bubble Tea app for the SSH session
func teaHandler(sshSession ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Stylizes the foreground (text) to be green
	lipgloss.SetColorProfile(termenv.TrueColor)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	tuiModel := tui.Model{
		Style: style,
	}

	// WithAltScreen makes it so that the app opens in another screen which preserves
	// your terminal history
	return tuiModel, []tea.ProgramOption{tea.WithAltScreen()}
}
