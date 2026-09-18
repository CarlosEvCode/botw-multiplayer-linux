package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Terminal Adaptive ANSI Colors (Inherits Kitty/Foot/Alacritty/Terminal Theme)
	ColorPrimary   = lipgloss.Color("4")  // Blue / Accent
	ColorHighlight = lipgloss.Color("3")  // Yellow / Gold
	ColorSuccess   = lipgloss.Color("2")  // Green
	ColorDanger    = lipgloss.Color("1")  // Red
	ColorWarning   = lipgloss.Color("3")  // Yellow/Orange
	ColorMuted     = lipgloss.Color("8")  // Bright Black / Gray
	ColorBorder    = lipgloss.Color("8")  // Border
	ColorFocus     = lipgloss.Color("6")  // Cyan / Active border
	ColorText      = lipgloss.Color("7")  // Foreground Text

	// Box & Title Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHighlight).
			Padding(0, 1)

	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(ColorPrimary).
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
			BorderForeground(ColorFocus).
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
			Foreground(ColorPrimary)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHighlight)

	DescStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
