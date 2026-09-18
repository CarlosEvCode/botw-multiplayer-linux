package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Palette
	ColorZeldaGreen = lipgloss.Color("#75B855")
	ColorSheikahBlue = lipgloss.Color("#3DB8FF")
	ColorTriforceGold = lipgloss.Color("#E6B800")
	ColorDarkBg     = lipgloss.Color("#1A1B26")
	ColorPanelBg    = lipgloss.Color("#16161E")
	ColorBorder     = lipgloss.Color("#3B4261")
	ColorActiveTab  = lipgloss.Color("#7AA2F7")
	ColorText       = lipgloss.Color("#C0CAF5")
	ColorMuted      = lipgloss.Color("#565F89")
	ColorSuccess    = lipgloss.Color("#9ECE6A")
	ColorDanger     = lipgloss.Color("#F7768E")
	ColorWarning    = lipgloss.Color("#E0AF68")

	// Box Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorTriforceGold).
			Padding(0, 1)

	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#1A1B26")).
			Background(ColorActiveTab).
			Padding(0, 2)

	TabInactiveStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 2)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	ActiveBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSheikahBlue).
			Padding(0, 1)

	// Status Badges
	StatusRunning = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess).
			SetString("[RUNNING]")

	StatusStopped = lipgloss.NewStyle().
			Foreground(ColorMuted).
			SetString("[STOPPED]")

	StatusActive = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			SetString("[ACTIVE]")

	StatusOk = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			SetString("[OK]")

	StatusWarn = lipgloss.NewStyle().
			Foreground(ColorWarning).
			SetString("[WARN]")

	StatusError = lipgloss.NewStyle().
			Foreground(ColorDanger).
			SetString("[ERROR]")

	// Text Highlights
	LabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSheikahBlue)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorTriforceGold)

	DescStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
