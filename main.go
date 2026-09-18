package main

import (
	"fmt"
	"os"

	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/config"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/process"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cargando configuracion: %v\n", err)
		os.Exit(1)
	}

	proc := process.NewProcessManager(500)
	model := ui.NewModel(cfg, proc)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error ejecutando interfaz TUI: %v\n", err)
		os.Exit(1)
	}

	// Clean up processes on exit if running
	if proc.IsServerRunning() {
		_ = proc.StopServer()
	}
}
