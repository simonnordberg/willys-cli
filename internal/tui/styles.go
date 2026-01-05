package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorPrimary   = lipgloss.Color("205")
	colorSecondary = lipgloss.Color("240")
	colorSuccess   = lipgloss.Color("green")
	colorDanger    = lipgloss.Color("red")
	colorMuted     = lipgloss.Color("241")
	colorStatusBg  = lipgloss.Color("236")
	colorStatusFg  = lipgloss.Color("229")

	statusBarStyle = lipgloss.NewStyle().
			Background(colorStatusBg).
			Foreground(colorStatusFg).
			Padding(0, 1)

	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Padding(0, 2)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(colorSecondary).
				Padding(0, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorDanger).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(colorSuccess)

	priceStyle = lipgloss.NewStyle().
			Foreground(colorSuccess)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)
)
