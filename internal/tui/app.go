package tui

import (
	"cloud-manage/internal/config"
	"cloud-manage/internal/consts"
	"cloud-manage/internal/tui/views"
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
	regionInput  textinput.Model
	focusIndex   int // 0=accessKeyId, 1=secret, 2=region

	// Credentials (set after login)
	ak     string
	sk     string
	region string

	// Views
	ecsView views.ECSView
	cmsView views.CMSView
	slsView views.SLSView
	ossView views.OSSView
	vpcView views.VPCView
	slbView views.SLBView

	// Track if current view has been loaded
	loaded map[Tab]bool
}

func NewApp() App {
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	showLogin := accessKeyId == "" || accessKeySecret == ""

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

	regionInput := textinput.New()
	regionInput.Placeholder = "cn-hangzhou"
	regionInput.CharLimit = 30
	regionInput.Width = 30
	regionInput.Prompt = "Region:        "
	regionInput.SetValue("cn-hangzhou")

	app := App{
		activeTab:    TabECS,
		showLogin:    showLogin,
		accessKeyId:  idInput,
		accessSecret: secretInput,
		regionInput:  regionInput,
		focusIndex:   0,
		region:       "cn-hangzhou",
		loaded:       make(map[Tab]bool),
	}

	// Load theme from config
	if cfg, err := config.Load(); err == nil && cfg.Theme != "" {
		SetTheme(cfg.Theme)
	}

	if showLogin {
		app.accessKeyId.Focus()
	}

	return app
}

func (a App) Init() tea.Cmd {
	return textinput.Blink
}

func (a *App) initViews() {
	a.ecsView = views.NewECSView()
	a.cmsView = views.NewCMSView()
	a.slsView = views.NewSLSView()
	a.ossView = views.NewOSSView()
	a.vpcView = views.NewVPCView()
	a.slbView = views.NewSLBView()

	a.ecsView.SetCredentials(a.ak, a.sk, a.region)
	a.cmsView.SetCredentials(a.ak, a.sk, a.region)
	a.slsView.SetCredentials(a.ak, a.sk, a.region)
	a.ossView.SetCredentials(a.ak, a.sk, a.region)
	a.vpcView.SetCredentials(a.ak, a.sk, a.region)
	a.slbView.SetCredentials(a.ak, a.sk, a.region)

	a.ecsView.SetSize(a.width, a.height-6)
	a.cmsView.SetSize(a.width, a.height-6)
	a.slsView.SetSize(a.width, a.height-6)
	a.ossView.SetSize(a.width, a.height-6)
	a.vpcView.SetSize(a.width, a.height-6)
	a.slbView.SetSize(a.width, a.height-6)
}

func (a App) loadCurrentView() tea.Cmd {
	switch a.activeTab {
	case TabECS:
		if !a.loaded[TabECS] {
			a.loaded[TabECS] = true
			return a.ecsView.LoadData(a.ak, a.sk, a.region)
		}
	case TabCMS:
		// CMS needs an instance ID, skip auto-load
	case TabSLS:
		// SLS needs project/logstore, skip auto-load
	case TabOSS:
		if !a.loaded[TabOSS] {
			a.loaded[TabOSS] = true
			return a.ossView.LoadBuckets(a.ak, a.sk, a.region)
		}
	case TabVPC:
		if !a.loaded[TabVPC] {
			a.loaded[TabVPC] = true
			return a.vpcView.LoadData(a.ak, a.sk, a.region)
		}
	case TabSLB:
		if !a.loaded[TabSLB] {
			a.loaded[TabSLB] = true
			return a.slbView.LoadData(a.ak, a.sk, a.region)
		}
	}
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
		a.ecsView.SetSize(a.width, a.height-6)
		a.cmsView.SetSize(a.width, a.height-6)
		a.slsView.SetSize(a.width, a.height-6)
		a.ossView.SetSize(a.width, a.height-6)
		a.vpcView.SetSize(a.width, a.height-6)
		a.slbView.SetSize(a.width, a.height-6)
		return a, nil

	case tea.KeyMsg:
		// Login mode
		if a.showLogin {
			return a.updateLogin(msg)
		}

		// Global keys
		switch {
		case key.Matches(msg, Keys.Quit):
			return a, tea.Quit
		case key.Matches(msg, Keys.Left):
			if a.activeTab > TabECS {
				a.activeTab--
				cmds = append(cmds, a.loadCurrentView())
			}
		case key.Matches(msg, Keys.Right):
			if a.activeTab < TabSLB {
				a.activeTab++
				cmds = append(cmds, a.loadCurrentView())
			}
		case key.Matches(msg, Keys.Tab1):
			a.activeTab = TabECS
			cmds = append(cmds, a.loadCurrentView())
		case key.Matches(msg, Keys.Tab2):
			a.activeTab = TabCMS
			cmds = append(cmds, a.loadCurrentView())
		case key.Matches(msg, Keys.Tab3):
			a.activeTab = TabSLS
			cmds = append(cmds, a.loadCurrentView())
		case key.Matches(msg, Keys.Tab4):
			a.activeTab = TabOSS
			cmds = append(cmds, a.loadCurrentView())
		case key.Matches(msg, Keys.Tab5):
			a.activeTab = TabVPC
			cmds = append(cmds, a.loadCurrentView())
		case key.Matches(msg, Keys.Tab6):
			a.activeTab = TabSLB
			cmds = append(cmds, a.loadCurrentView())
		case key.Matches(msg, Keys.Theme):
			// Cycle theme: dark -> light -> dark
			if GetTheme() == "dark" {
				SetTheme("light")
			} else {
				SetTheme("dark")
			}
			config.UpdateTheme(GetTheme())
		}

		// Route to active view
		cmd := a.routeToView(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return a, tea.Batch(cmds...)

	default:
		// Route other messages (async results) to active view
		cmd := a.routeToView(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return a, tea.Batch(cmds...)
	}
}

func (a *App) updateLogin(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.String() {
	case "tab", "down":
		a.focusIndex = (a.focusIndex + 1) % 3
		a.updateFocus()
	case "shift+tab", "up":
		a.focusIndex = (a.focusIndex + 2) % 3
		a.updateFocus()
	case "enter":
		if a.focusIndex == 2 {
			// Submit
			a.ak = a.accessKeyId.Value()
			a.sk = a.accessSecret.Value()
			if r := a.regionInput.Value(); r != "" {
				a.region = r
			}
			a.showLogin = false
			os.Setenv("CLOUD_ACCESS_KEY_ID", a.ak)
			os.Setenv("CLOUD_ACCESS_KEY_SECRET", a.sk)
			a.initViews()
			cmds = append(cmds, a.loadCurrentView())
			return a, tea.Batch(cmds...)
		}
		a.focusIndex++
		a.updateFocus()
	case "esc":
		return a, tea.Quit
	default:
		if a.focusIndex == 0 {
			a.accessKeyId, cmd = a.accessKeyId.Update(msg)
		} else if a.focusIndex == 1 {
			a.accessSecret, cmd = a.accessSecret.Update(msg)
		} else {
			a.regionInput, cmd = a.regionInput.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) updateFocus() {
	a.accessKeyId.Blur()
	a.accessSecret.Blur()
	a.regionInput.Blur()

	switch a.focusIndex {
	case 0:
		a.accessKeyId.Focus()
	case 1:
		a.accessSecret.Focus()
	case 2:
		a.regionInput.Focus()
	}
}

func (a *App) routeToView(msg tea.Msg) tea.Cmd {
	switch a.activeTab {
	case TabECS:
		return a.ecsView.HandleMessage(msg)
	case TabCMS:
		return a.cmsView.HandleMessage(msg)
	case TabSLS:
		return a.slsView.HandleMessage(msg)
	case TabOSS:
		return a.ossView.HandleMessage(msg)
	case TabVPC:
		return a.vpcView.HandleMessage(msg)
	case TabSLB:
		return a.slbView.HandleMessage(msg)
	}
	return nil
}

func (a App) View() string {
	if !a.ready {
		return "Initializing..."
	}
	if a.showLogin {
		return a.renderLogin()
	}
	return a.renderMain()
}

func (a App) renderLogin() string {
	title := TitleBarStyle.Width(a.width).Render("Cloud Manage - 登录")

	content := fmt.Sprintf(`
  请输入阿里云凭证：

  %s
  %s
  %s

  按 Tab/↑↓ 切换输入框，Enter 确认，Esc 退出
`, a.accessKeyId.View(), a.accessSecret.View(), a.regionInput.View())

	return fmt.Sprintf("%s\n%s", title, content)
}

func (a App) renderMain() string {
	title := TitleBarStyle.Width(a.width).Render(fmt.Sprintf("Cloud Manage TUI %s | Region: %s", consts.Version, a.region))
	tabs := renderTabs(a.activeTab)

	var content string
	switch a.activeTab {
	case TabECS:
		content = a.ecsView.Render()
	case TabCMS:
		content = a.cmsView.Render()
	case TabSLS:
		content = a.slsView.Render()
	case TabOSS:
		content = a.ossView.Render()
	case TabVPC:
		content = a.vpcView.Render()
	case TabSLB:
		content = a.slbView.Render()
	}

	status := StatusBarStyle.Width(a.width).Render("  1-6 Switch Tab | ←→ Tab | ↑↓/jk Navigate | Enter Detail | Esc Back | q Quit")

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
