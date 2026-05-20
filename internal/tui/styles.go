package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Primary   = lipgloss.Color("#00BFFF")
	Secondary = lipgloss.Color("#FF6B6B")
	Success   = lipgloss.Color("#00FF88")
	Warning   = lipgloss.Color("#FFD700")
	Error     = lipgloss.Color("#FF4444")
	Gray      = lipgloss.Color("#666666")
	White     = lipgloss.Color("#FFFFFF")

	// Title bar
	TitleBarStyle = lipgloss.NewStyle().
			Background(Primary).
			Foreground(White).
			Bold(true).
			Padding(0, 1)

	// Tabs
	TabStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(Gray)

	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(Primary)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(White).
			Padding(0, 1)
)
