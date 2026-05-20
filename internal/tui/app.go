package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
)

type Tab int

const (
	TabECS Tab = iota
	TabCMS
	TabSLS
	TabOSS
)

var tabNames = []string{"ECS", "CMS", "SLS", "OSS"}

type App struct {
	activeTab Tab
	width     int
	height    int
	ready     bool
}

func NewApp() App {
	return App{
		activeTab: TabECS,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Quit):
			return a, tea.Quit
		case key.Matches(msg, Keys.Left):
			if a.activeTab > TabECS {
				a.activeTab--
			}
		case key.Matches(msg, Keys.Right):
			if a.activeTab < TabOSS {
				a.activeTab++
			}
		}
	}
	return a, nil
}

func (a App) View() string {
	if !a.ready {
		return "Initializing..."
	}

	// Title bar
	title := TitleBarStyle.Width(a.width).Render("Cloud Manage TUI v0.0.13")

	// Tabs
	tabs := renderTabs()

	// Content
	content := fmt.Sprintf("\n  Content for %s tab\n", tabNames[a.activeTab])

	// Status
	status := StatusBarStyle.Width(a.width).Render("  ↑↓ Navigate | ←→ Tab | Enter Select | q Quit")

	return fmt.Sprintf("%s\n%s\n%s\n%s", title, tabs, content, status)
}

func renderTabs() string {
	var rendered string
	for i, name := range tabNames {
		if i == int(0) { // TODO: use activeTab
			rendered += ActiveTabStyle.Render(name) + " "
		} else {
			rendered += TabStyle.Render(name) + " "
		}
	}
	return rendered
}
