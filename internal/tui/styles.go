package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Current theme: "dark" or "light"
	currentTheme = "dark"

	// Colors (updated by theme switching)
	Primary   = lipgloss.Color("#00BFFF")
	Secondary = lipgloss.Color("#FF6B6B")
	Success   = lipgloss.Color("#00FF88")
	Warning   = lipgloss.Color("#FFD700")
	Error     = lipgloss.Color("#FF4444")
	Gray      = lipgloss.Color("#666666")
	White     = lipgloss.Color("#FFFFFF")
	BGColor   = lipgloss.Color("#1A1A2E")
	TextColor = lipgloss.Color("#E0E0E0")

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

// SetTheme switches between dark and light themes.
func SetTheme(theme string) {
	currentTheme = theme
	switch theme {
	case "light":
		Primary = lipgloss.Color("#0066CC")
		Secondary = lipgloss.Color("#CC3333")
		Success = lipgloss.Color("#228B22")
		Warning = lipgloss.Color("#CC8800")
		Error = lipgloss.Color("#CC0000")
		Gray = lipgloss.Color("#999999")
		White = lipgloss.Color("#333333")
		BGColor = lipgloss.Color("#F5F5F5")
		TextColor = lipgloss.Color("#333333")

		TitleBarStyle = lipgloss.NewStyle().
				Background(Primary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true).
				Padding(0, 1)

		TabStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(Gray)

		ActiveTabStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Primary).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(Primary)

		StatusBarStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#DDDDDD")).
				Foreground(lipgloss.Color("#333333")).
				Padding(0, 1)

	default: // dark
		Primary = lipgloss.Color("#00BFFF")
		Secondary = lipgloss.Color("#FF6B6B")
		Success = lipgloss.Color("#00FF88")
		Warning = lipgloss.Color("#FFD700")
		Error = lipgloss.Color("#FF4444")
		Gray = lipgloss.Color("#666666")
		White = lipgloss.Color("#FFFFFF")
		BGColor = lipgloss.Color("#1A1A2E")
		TextColor = lipgloss.Color("#E0E0E0")

		TitleBarStyle = lipgloss.NewStyle().
				Background(Primary).
				Foreground(White).
				Bold(true).
				Padding(0, 1)

		TabStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(Gray)

		ActiveTabStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Primary).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(Primary)

		StatusBarStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#333333")).
				Foreground(White).
				Padding(0, 1)
	}
}

// GetTheme returns the current theme name.
func GetTheme() string {
	return currentTheme
}
