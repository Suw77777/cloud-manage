package tui

import (
	"cloud-manage/internal/consts"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Tab int

const (
	TabECS Tab = iota
	TabCMS
	TabSLS
	TabOSS
	TabVPC
	TabSLB
)

var tabNames = []string{"ECS", "CMS", "SLS", "OSS", "VPC", "SLB"}

type App struct {
	activeTab    Tab
	width        int
	height       int
	ready        bool
	showLogin    bool
	accessKeyId  textinput.Model
	accessSecret textinput.Model
	focusIndex   int // 0=accessKeyId, 1=secret
}

func NewApp() App {
	// Check if credentials exist
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")

	showLogin := accessKeyId == "" || accessKeySecret == ""

	// Create input fields
	idInput := textinput.New()
	idInput.Placeholder = "AccessKey ID"
	idInput.CharLimit = 50
	idInput.Width = 50
	idInput.Prompt = "AccessKey ID: "

	secretInput := textinput.New()
	secretInput.Placeholder = "AccessKey Secret"
	secretInput.CharLimit = 100
	secretInput.Width = 50
	secretInput.EchoMode = textinput.EchoPassword
	secretInput.Prompt = "AccessKey Secret: "

	app := App{
		activeTab:    TabECS,
		showLogin:    showLogin,
		accessKeyId:  idInput,
		accessSecret: secretInput,
		focusIndex:   0,
	}

	if showLogin {
		app.accessKeyId.Focus()
	}

	return app
}

func (a App) Init() tea.Cmd {
	return textinput.Blink
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true

	case tea.KeyMsg:
		// Login mode
		if a.showLogin {
			switch msg.String() {
			case "tab", "enter":
				if a.focusIndex == 0 {
					a.focusIndex = 1
					a.accessKeyId.Blur()
					a.accessSecret.Focus()
				} else {
					// Submit credentials
					a.showLogin = false
					os.Setenv("CLOUD_ACCESS_KEY_ID", a.accessKeyId.Value())
					os.Setenv("CLOUD_ACCESS_KEY_SECRET", a.accessSecret.Value())
				}
			case "esc":
				return a, tea.Quit
			default:
				if a.focusIndex == 0 {
					a.accessKeyId, cmd = a.accessKeyId.Update(msg)
				} else {
					a.accessSecret, cmd = a.accessSecret.Update(msg)
				}
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}

		// Normal mode
		switch {
		case key.Matches(msg, Keys.Quit):
			return a, tea.Quit
		case key.Matches(msg, Keys.Left):
			if a.activeTab > TabECS {
				a.activeTab--
			}
		case key.Matches(msg, Keys.Right):
			if a.activeTab < TabSLB {
				a.activeTab++
			}
		case key.Matches(msg, Keys.Tab1):
			a.activeTab = TabECS
		case key.Matches(msg, Keys.Tab2):
			a.activeTab = TabCMS
		case key.Matches(msg, Keys.Tab3):
			a.activeTab = TabSLS
		case key.Matches(msg, Keys.Tab4):
			a.activeTab = TabOSS
		case key.Matches(msg, Keys.Tab5):
			a.activeTab = TabVPC
		case key.Matches(msg, Keys.Tab6):
			a.activeTab = TabSLB
		}
	}

	// Update text inputs
	if a.showLogin {
		if a.focusIndex == 0 {
			a.accessKeyId, cmd = a.accessKeyId.Update(msg)
		} else {
			a.accessSecret, cmd = a.accessSecret.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if !a.ready {
		return "Initializing..."
	}

	// Login screen
	if a.showLogin {
		return a.renderLogin()
	}

	// Main screen
	return a.renderMain()
}

func (a App) renderLogin() string {
	title := TitleBarStyle.Width(a.width).Render("Cloud Manage - 登录")

	content := fmt.Sprintf(`
  请输入阿里云凭证：

  %s
  %s

  按 Tab 切换输入框，Enter 确认，Esc 退出
`, a.accessKeyId.View(), a.accessSecret.View())

	return fmt.Sprintf("%s\n%s", title, content)
}

func (a App) renderMain() string {
	// Title bar
	title := TitleBarStyle.Width(a.width).Render("Cloud Manage TUI " + consts.Version)

	// Tabs
	tabs := renderTabs(a.activeTab)

	// Content
	content := fmt.Sprintf("\n  [%s] 内容区域\n\n  按 ←→ 切换 Tab，q 退出", tabNames[a.activeTab])

	// Status
	status := StatusBarStyle.Width(a.width).Render("  ↑↓ Navigate | ←→ Tab | Enter Select | q Quit")

	return fmt.Sprintf("%s\n%s\n\n%s\n%s", title, tabs, content, status)
}

func renderTabs(active Tab) string {
	var rendered string
	for i, name := range tabNames {
		if i == int(active) {
			rendered += ActiveTabStyle.Render(name) + " "
		} else {
			rendered += TabStyle.Render(name) + " "
		}
	}
	return rendered
}
