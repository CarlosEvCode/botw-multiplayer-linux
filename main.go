package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/config"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/process"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/ui"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/uninstall"
	tea "github.com/charmbracelet/bubbletea"
)

const version = "v1.0.0"

func main() {
	uninstallFlag := flag.Bool("uninstall", false, "Completely uninstall and remove the Wineprefix, configs, and launcher files")
	uninstallShort := flag.Bool("u", false, "Alias for --uninstall")
	yesFlag := flag.Bool("yes", false, "Automatically confirm uninstallation without prompting")
	yesShort := flag.Bool("y", false, "Alias for --yes")
	versionFlag := flag.Bool("version", false, "Display application version")
	versionShort := flag.Bool("v", false, "Alias for --version")

	flag.Usage = func() {
		fmt.Printf("Zelda: Breath of the Wild Multiplayer Manager (%s)\n\n", version)
		fmt.Println("Usage:")
		fmt.Println("  botw-manager [flags]")
		fmt.Println("\nFlags:")
		fmt.Println("  -u, --uninstall   Completely remove Wineprefix, launchers, configs, and binaries")
		fmt.Println("  -y, --yes         Auto-confirm prompt when running uninstaller")
		fmt.Println("  -v, --version     Display application version")
		fmt.Println("  -h, --help        Display this help message")
	}

	flag.Parse()

	if *versionFlag || *versionShort {
		fmt.Printf("botw-manager %s\n", version)
		os.Exit(0)
	}

	if *uninstallFlag || *uninstallShort {
		autoConfirm := *yesFlag || *yesShort
		if err := uninstall.RunUninstaller(autoConfirm); err != nil {
			fmt.Fprintf(os.Stderr, "Error during uninstallation: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

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
