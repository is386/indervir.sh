package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/is386/indervir.dev/src/tui"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/ssh"
)

const (
	host = "0.0.0.0"
	port = 2235
)

func main() {
	// Force TrueColor so styles (bold, colors) work over SSH
	lipgloss.Writer.Profile = colorprofile.TrueColor

	// Generates a new SSH server with the given address, host key, and middleware
	// First, activeterm rejects connections without a PTY (terminal)
	// Second, bubbletea launches a Bubble Tea TUI app for the SSH session
	wishServer, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				return tui.InitialModel(), nil
			}),
			activeterm.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		os.Exit(1)
	}

	// Starts the server in a goroutine and listens for SIGINT/SIGTERM (ctrl+c)
	serverDone := make(chan os.Signal, 1)
	signal.Notify(serverDone, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err := wishServer.ListenAndServe(); err != nil {
			log.Error("Could not start server", "error", err)
			serverDone <- nil
		}
	}()

	// Graceful shutdown
	<-serverDone
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := wishServer.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}
