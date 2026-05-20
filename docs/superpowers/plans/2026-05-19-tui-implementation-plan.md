# TUI 终端图形界面实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 cloud-manage 添加基于 Bubble Tea 的终端图形界面（TUI），支持 ECS/CMS/SLS/OSS 完整功能，实时刷新、搜索过滤、多 region 切换。

**Architecture:** 使用 Bubble Tea 框架实现 Elm 架构（Model-Update-View），复用现有 service 层，TUI 作为新的展示层。

**Tech Stack:** Go, Bubble Tea, Bubbles, Lip Gloss

---

## 文件结构

### 新建文件

| 文件 | 职责 |
|------|------|
| `cmd/tui/main.go` | TUI 入口 |
| `internal/tui/app.go` | 主应用模型 |
| `internal/tui/styles.go` | 样式定义 |
| `internal/tui/keys.go` | 快捷键定义 |
| `internal/tui/views/ecs.go` | ECS 视图 |
| `internal/tui/views/cms.go` | CMS 视图 |
| `internal/tui/views/sls.go` | SLS 视图 |
| `internal/tui/views/oss.go` | OSS 视图 |
| `internal/tui/components/table.go` | 表格组件 |
| `internal/tui/components/detail.go` | 详情面板 |
| `internal/tui/components/search.go` | 搜索框 |
| `internal/tui/components/status.go` | 状态栏 |
| `internal/tui/components/tabs.go` | Tab 切换 |

### 修改文件

| 文件 | 修改内容 |
|------|----------|
| `main.go` | 添加 TUI 模式检测 |
| `go.mod` | 添加 Bubble Tea 依赖 |

---

## Task 1: 安装依赖并创建基础结构

**Files:**
- Modify: `go.mod`
- Create: `cmd/tui/main.go`
- Create: `internal/tui/app.go`

- [ ] **Step 1: 安装 Bubble Tea 依赖**

```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/charmbracelet/lipgloss@latest
```

- [ ] **Step 2: 创建 TUI 入口**

```go
// cmd/tui/main.go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"cloud-manage/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.NewApp(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: 创建主应用模型**

```go
// internal/tui/app.go
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Tab int

const (
	TabECS Tab = iota
	TabCMS
	TabSLS
	TabOSS
)

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
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "right", "l":
			if a.activeTab < TabOSS {
				a.activeTab++
			}
		case "left", "h":
			if a.activeTab > TabECS {
				a.activeTab--
			}
		}
	}
	return a, nil
}

func (a App) View() string {
	if !a.ready {
		return "Initializing..."
	}
	return "Cloud Manage TUI"
}
```

- [ ] **Step 4: 验证编译**

```bash
go build ./cmd/tui/
```

- [ ] **Step 5: 提交**

```bash
git add go.mod go.sum cmd/tui/ internal/tui/
git commit -m "feat: add TUI basic structure with Bubble Tea"
```

---

## Task 2: 创建样式和快捷键定义

**Files:**
- Create: `internal/tui/styles.go`
- Create: `internal/tui/keys.go`

- [ ] **Step 1: 创建样式定义**

```go
// internal/tui/styles.go
package tui

import "github.com/charmbracelet/lipgloss"

var (
	// 颜色
	Primary   = lipgloss.Color("#00BFFF")
	Secondary = lipgloss.Color("#FF6B6B")
	Success   = lipgloss.Color("#00FF88")
	Warning   = lipgloss.Color("#FFD700")
	Error     = lipgloss.Color("#FF4444")
	Gray      = lipgloss.Color("#666666")
	White     = lipgloss.Color("#FFFFFF")

	// 标题栏
	TitleBarStyle = lipgloss.NewStyle().
			Background(Primary).
			Foreground(White).
			Bold(true).
			Padding(0, 1)

	// Tab 栏
	TabStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(Gray)

	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(Primary)

	// 表格
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Primary).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(Gray)

	TableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	SelectedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1A1A2E")).
				Foreground(White)

	// 状态栏
	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(White).
			Padding(0, 1)

	// 详情面板
	DetailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	// 搜索框
	SearchStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)
)
```

- [ ] **Step 2: 创建快捷键定义**

```go
// internal/tui/keys.go
package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Search   key.Binding
	Refresh  key.Binding
	Region   key.Binding
	Quit     key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev tab"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next tab"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Region: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "change region"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
```

- [ ] **Step 3: 提交**

```bash
git add internal/tui/styles.go internal/tui/keys.go
git commit -m "feat: add TUI styles and key bindings"
```

---

## Task 3: 创建组件 - Tab 切换

**Files:**
- Create: `internal/tui/components/tabs.go`

- [ ] **Step 1: 创建 Tab 组件**

```go
// internal/tui/components/tabs.go
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	tabBorder = lipgloss.NormalBorder()

	tabStyle = lipgloss.NewStyle().
			Border(tabBorder, false, false, true, false).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(0, 1)

	activeTabStyle = tabStyle.Copy().
			BorderForeground(lipgloss.Color("#00BFFF")).
			Foreground(lipgloss.Color("#00BFFF")).
			Bold(true)
)

func RenderTabs(tabs []string, active int, width int) string {
	var rendered []string
	for i, tab := range tabs {
		if i == active {
			rendered = append(rendered, activeTabStyle.Render(tab))
		} else {
			rendered = append(rendered, tabStyle.Render(tab))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	return lipgloss.Place(width, 1, lipgloss.Left, lipgloss.Center, row)
}
```

- [ ] **Step 2: 提交**

```bash
git add internal/tui/components/tabs.go
git commit -m "feat: add tabs component"
```

---

## Task 4: 创建组件 - 表格

**Files:**
- Create: `internal/tui/components/table.go`

- [ ] **Step 1: 创建表格组件**

```go
// internal/tui/components/table.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Table struct {
	Headers  []string
	Rows     [][]string
	Cursor   int
	Offset   int
	MaxRows  int
	width    int
	height   int
}

func NewTable(headers []string) Table {
	return Table{
		Headers: headers,
		MaxRows: 100,
	}
}

func (t *Table) SetSize(width, height int) {
	t.width = width
	t.height = height
}

func (t *Table) SetRows(rows [][]string) {
	t.Rows = rows
	if t.Cursor >= len(rows) {
		t.Cursor = len(rows) - 1
	}
	if t.Cursor < 0 {
		t.Cursor = 0
	}
}

func (t *Table) MoveUp() {
	if t.Cursor > 0 {
		t.Cursor--
		if t.Cursor < t.Offset {
			t.Offset = t.Cursor
		}
	}
}

func (t *Table) MoveDown() {
	if t.Cursor < len(t.Rows)-1 {
		t.Cursor++
		visibleRows := t.height - 3
		if t.Cursor >= t.Offset+visibleRows {
			t.Offset = t.Cursor - visibleRows + 1
		}
	}
}

func (t *Table) SelectedRow() []string {
	if t.Cursor >= 0 && t.Cursor < len(t.Rows) {
		return t.Rows[t.Cursor]
	}
	return nil
}

func (t Table) Render() string {
	if len(t.Headers) == 0 {
		return ""
	}

	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00BFFF")).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#666666"))

	for i, h := range t.Headers {
		if i > 0 {
			b.WriteString("  ")
		}
		b.WriteString(lipgloss.NewStyle().Width(20).Render(h))
	}
	b.WriteString("\n")
	b.WriteString(headerStyle.Render(b.String()))
	b.Reset()

	// Rows
	visibleRows := t.height - 3
	if visibleRows < 1 {
		visibleRows = 1
	}

	end := t.Offset + visibleRows
	if end > len(t.Rows) {
		end = len(t.Rows)
	}

	for i := t.Offset; i < end; i++ {
		row := t.Rows[i]
		rowStr := ""
		for j, cell := range row {
			if j > 0 {
				rowStr += "  "
			}
			rowStr += fmt.Sprintf("%-20s", cell)
		}

		if i == t.Cursor {
			selectedStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("#1A1A2E")).
				Foreground(lipgloss.Color("#FFFFFF"))
			b.WriteString(selectedStyle.Render(rowStr))
		} else {
			b.WriteString(rowStr)
		}
		b.WriteString("\n")
	}

	// Padding if needed
	for i := end - t.Offset; i < visibleRows; i++ {
		b.WriteString("\n")
	}

	return b.String()
}
```

- [ ] **Step 2: 提交**

```bash
git add internal/tui/components/table.go
git commit -m "feat: add table component"
```

---

## Task 5: 创建组件 - 详情面板

**Files:**
- Create: `internal/tui/components/detail.go`

- [ ] **Step 1: 创建详情面板**

```go
// internal/tui/components/detail.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Detail struct {
	Title   string
	Fields  map[string]string
	width   int
	height  int
}

func NewDetail(title string) Detail {
	return Detail{
		Title:  title,
		Fields: make(map[string]string),
	}
}

func (d *Detail) SetSize(width, height int) {
	d.width = width
	d.height = height
}

func (d *Detail) SetField(key, value string) {
	d.Fields[key] = value
}

func (d *Detail) Clear() {
	d.Fields = make(map[string]string)
}

func (d Detail) Render() string {
	if len(d.Fields) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("No details available")
	}

	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00BFFF")).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#666666")).
		Width(d.width - 4)

	b.WriteString(titleStyle.Render(d.Title))
	b.WriteString("\n\n")

	// Fields
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF")).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	for key, value := range d.Fields {
		b.WriteString(labelStyle.Render(key + ":"))
		b.WriteString(valueStyle.Render(value))
		b.WriteString("\n")
	}

	// Border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00BFFF")).
		Padding(1, 2).
		Width(d.width - 2)

	return borderStyle.Render(b.String())
}
```

- [ ] **Step 2: 提交**

```bash
git add internal/tui/components/detail.go
git commit -m "feat: add detail panel component"
```

---

## Task 6: 创建组件 - 搜索框和状态栏

**Files:**
- Create: `internal/tui/components/search.go`
- Create: `internal/tui/components/status.go`

- [ ] **Step 1: 创建搜索框**

```go
// internal/tui/components/search.go
package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Search struct {
	input    textinput.Model
	active   bool
	width    int
}

func NewSearch() Search {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 100
	ti.Width = 50
	return Search{
		input: ti,
	}
}

func (s *Search) SetWidth(width int) {
	s.width = width
	s.input.Width = width - 4
}

func (s *Search) Activate() {
	s.active = true
	s.input.Focus()
}

func (s *Search) Deactivate() {
	s.active = false
	s.input.Blur()
	s.input.SetValue("")
}

func (s *Search) IsActive() bool {
	return s.active
}

func (s *Search) Value() string {
	return s.input.Value()
}

func (s *Search) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return cmd
}

func (s Search) Render() string {
	if !s.active {
		return ""
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00BFFF")).
		Padding(0, 1).
		Width(s.width - 2)

	return style.Render(s.input.View())
}
```

- [ ] **Step 2: 创建状态栏**

```go
// internal/tui/components/status.go
package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Status struct {
	Message string
	Error   string
	Region  string
	width   int
}

func NewStatus() Status {
	return Status{
		Region: "cn-hangzhou",
	}
}

func (s *Status) SetWidth(width int) {
	s.width = width
}

func (s *Status) SetMessage(msg string) {
	s.Message = msg
	s.Error = ""
}

func (s *Status) SetError(err string) {
	s.Error = err
	s.Message = ""
}

func (s *Status) SetRegion(region string) {
	s.Region = region
}

func (s Status) Render() string {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Width(s.width)

	left := ""
	if s.Error != "" {
		left = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Render("Error: " + s.Error)
	} else if s.Message != "" {
		left = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF88")).
			Render(s.Message)
	} else {
		left = "Ready"
	}

	right := fmt.Sprintf("Region: %s | q: quit | /: search | r: refresh", s.Region)

	gap := s.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 0 {
		gap = 0
	}

	return style.Render(left + lipgloss.NewStyle().Width(gap).Render("") + right)
}
```

- [ ] **Step 3: 提交**

```bash
git add internal/tui/components/search.go internal/tui/components/status.go
git commit -m "feat: add search and status components"
```

---

## Task 7: 创建 ECS 视图

**Files:**
- Create: `internal/tui/views/ecs.go`

- [ ] **Step 1: 创建 ECS 视图**

```go
// internal/tui/views/ecs.go
package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"cloud-manage/internal/tui/components"
	"cloud-manage/service"
)

type ECSView struct {
	table    components.Table
	detail   components.Detail
	loading  bool
	showDetail bool
	service  *service.ECSService
	width    int
	height   int
}

func NewECSView() ECSView {
	return ECSView{
		table:   components.NewTable([]string{"ID", "Name", "Status", "Public IP", "Private IP"}),
		detail:  components.NewDetail("Instance Detail"),
		service: service.NewECSService(),
	}
}

func (v *ECSView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.table.SetSize(width, height/2)
	v.detail.SetSize(width, height/2)
}

func (v *ECSView) LoadData(accessKeyId, accessKeySecret, region string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.service.ListInstances(accessKeyId, accessKeySecret, region)
		if err != nil {
			return ECSLoadError{err}
		}
		return ECSLoaded{result}
	}
}

type ECSLoaded struct {
	Result *service.ListInstancesResult
}

type ECSLoadError struct {
	Err error
}

func (v *ECSView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ECSLoaded:
		v.loading = false
		rows := make([][]string, 0)
		for _, inst := range msg.Result.Instances {
			rows = append(rows, []string{
				inst.InstanceId,
				inst.InstanceName,
				inst.Status,
				inst.PublicIp,
				inst.PrivateIp,
			})
		}
		v.table.SetRows(rows)
	case ECSLoadError:
		v.loading = false
	}
	return nil
}

func (v *ECSView) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		v.table.MoveUp()
	case "down", "j":
		v.table.MoveDown()
	case "enter":
		row := v.table.SelectedRow()
		if row != nil {
			v.showDetail = true
			v.detail.Title = row[0]
			v.detail.SetField("ID", row[0])
			v.detail.SetField("Name", row[1])
			v.detail.SetField("Status", row[2])
			v.detail.SetField("Public IP", row[3])
			v.detail.SetField("Private IP", row[4])
		}
	case "esc":
		v.showDetail = false
		v.detail.Clear()
	}
	return nil
}

func (v ECSView) Render() string {
	if v.loading {
		return "Loading..."
	}

	if v.showDetail {
		return v.detail.Render()
	}

	return v.table.Render()
}
```

- [ ] **Step 2: 提交**

```bash
git add internal/tui/views/ecs.go
git commit -m "feat: add ECS view"
```

---

## Task 8: 创建 CMS/SLS/OSS 视图

**Files:**
- Create: `internal/tui/views/cms.go`
- Create: `internal/tui/views/sls.go`
- Create: `internal/tui/views/oss.go`

- [ ] **Step 1: 创建 CMS 视图**

```go
// internal/tui/views/cms.go
package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"cloud-manage/internal/tui/components"
	"cloud-manage/service"
)

type CMSView struct {
	table    components.Table
	detail   components.Detail
	loading  bool
	service  *service.CMSService
	width    int
	height   int
}

func NewCMSView() CMSView {
	return CMSView{
		table:   components.NewTable([]string{"Metric", "Value", "Unit"}),
		detail:  components.NewDetail("Metric Detail"),
		service: service.NewCMSService(),
	}
}

func (v *CMSView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.table.SetSize(width, height/2)
	v.detail.SetSize(width, height/2)
}

func (v *CMSView) LoadData(accessKeyId, accessKeySecret, region, instanceId string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.service.GetInstanceMetrics(accessKeyId, accessKeySecret, region, instanceId)
		if err != nil {
			return CMSLoadError{err}
		}
		return CMSLoaded{result}
	}
}

type CMSLoaded struct {
	Result *service.ECSMetricAdapter
}

type CMSLoadError struct {
	Err error
}

func (v *CMSView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case CMSLoaded:
		v.loading = false
		rows := make([][]string, 0)
		if msg.Result.CPUUtilization != nil {
			rows = append(rows, []string{"CPU", fmt.Sprintf("%.2f", *msg.Result.CPUUtilization), "%"})
		}
		if msg.Result.MemoryUtilization != nil {
			rows = append(rows, []string{"Memory", fmt.Sprintf("%.2f", *msg.Result.MemoryUtilization), "%"})
		}
		v.table.SetRows(rows)
	case CMSLoadError:
		v.loading = false
	}
	return nil
}

func (v *CMSView) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		v.table.MoveUp()
	case "down", "j":
		v.table.MoveDown()
	}
	return nil
}

func (v CMSView) Render() string {
	if v.loading {
		return "Loading..."
	}
	return v.table.Render()
}
```

- [ ] **Step 2: 创建 SLS 视图**

```go
// internal/tui/views/sls.go
package views

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"cloud-manage/internal/tui/components"
	"cloud-manage/service"
)

type SLSView struct {
	table    components.Table
	detail   components.Detail
	loading  bool
	service  *service.SLSService
	width    int
	height   int
}

func NewSLSView() SLSView {
	return SLSView{
		table:   components.NewTable([]string{"Timestamp", "Level", "Message"}),
		detail:  components.NewDetail("Log Detail"),
		service: service.NewSLSService(),
	}
}

func (v *SLSView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.table.SetSize(width, height/2)
	v.detail.SetSize(width, height/2)
}

func (v *SLSView) LoadData(accessKeyId, accessKeySecret, region, project, logstore, query string) tea.Cmd {
	return func() tea.Msg {
		from := time.Now().Add(-1 * time.Hour).Unix()
		to := time.Now().Unix()
		result, err := v.service.QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query, from, to, 100)
		if err != nil {
			return SLSLoadError{err}
		}
		return SLSLoaded{result}
	}
}

type SLSLoaded struct {
	Result *service.LogQueryResult
}

type SLSLoadError struct {
	Err error
}

func (v *SLSView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case SLSLoaded:
		v.loading = false
		rows := make([][]string, 0)
		for _, entry := range msg.Result.Entries {
			ts := time.Unix(entry.Timestamp/1000, 0).Format("15:04:05")
			level := entry.Content["level"]
			if level == "" {
				level = "INFO"
			}
			message := entry.Content["message"]
			if message == "" {
				for k, val := range entry.Content {
					if k != "level" && k != "timestamp" {
						message = fmt.Sprintf("%s=%s", k, val)
						break
					}
				}
			}
			rows = append(rows, []string{ts, level, message})
		}
		v.table.SetRows(rows)
	case SLSLoadError:
		v.loading = false
	}
	return nil
}

func (v *SLSView) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		v.table.MoveUp()
	case "down", "j":
		v.table.MoveDown()
	case "enter":
		// Show log detail
	}
	return nil
}

func (v SLSView) Render() string {
	if v.loading {
		return "Loading..."
	}
	return v.table.Render()
}
```

- [ ] **Step 3: 创建 OSS 视图**

```go
// internal/tui/views/oss.go
package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"cloud-manage/internal/tui/components"
	"cloud-manage/service"
)

type OSSView struct {
	table      components.Table
	detail     components.Detail
	loading    bool
	service    *service.OSSService
	showBucket bool
	bucket     string
	width      int
	height     int
}

func NewOSSView() OSSView {
	return OSSView{
		table:   components.NewTable([]string{"Name", "Location", "Created"}),
		detail:  components.NewDetail("Bucket Detail"),
		service: service.NewOSSService(),
	}
}

func (v *OSSView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.table.SetSize(width, height/2)
	v.detail.SetSize(width, height/2)
}

func (v *OSSView) LoadBuckets(accessKeyId, accessKeySecret, region string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.service.ListBuckets(accessKeyId, accessKeySecret, region)
		if err != nil {
			return OSSLoadError{err}
		}
		return OSSBucketsLoaded{result}
	}
}

func (v *OSSView) LoadObjects(accessKeyId, accessKeySecret, region, bucket string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.service.ListObjects(accessKeyId, accessKeySecret, region, bucket, "", 100)
		if err != nil {
			return OSSLoadError{err}
		}
		return OSSObjectsLoaded{result}
	}
}

type OSSBucketsLoaded struct {
	Result *service.ListBucketsResult
}

type OSSObjectsLoaded struct {
	Result *service.ListObjectsResult
}

type OSSLoadError struct {
	Err error
}

func (v *OSSView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case OSSBucketsLoaded:
		v.loading = false
		v.showBucket = false
		rows := make([][]string, 0)
		for _, b := range msg.Result.Buckets {
			rows = append(rows, []string{b.Name, b.Location, b.CreationDate})
		}
		v.table = components.NewTable([]string{"Name", "Location", "Created"})
		v.table.SetSize(v.width, v.height/2)
		v.table.SetRows(rows)
	case OSSObjectsLoaded:
		v.loading = false
		v.showBucket = true
		rows := make([][]string, 0)
		for _, obj := range msg.Result.Objects {
			rows = append(rows, []string{obj.Key, formatSize(obj.Size), obj.LastModified})
		}
		v.table = components.NewTable([]string{"Name", "Size", "Modified"})
		v.table.SetSize(v.width, v.height/2)
		v.table.SetRows(rows)
	case OSSLoadError:
		v.loading = false
	}
	return nil
}

func (v *OSSView) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		v.table.MoveUp()
	case "down", "j":
		v.table.MoveDown()
	case "enter":
		if !v.showBucket {
			row := v.table.SelectedRow()
			if row != nil {
				v.bucket = row[0]
				return v.LoadObjects("", "", "", v.bucket)
			}
		}
	case "esc":
		if v.showBucket {
			v.showBucket = false
			return v.LoadBuckets("", "", "")
		}
	}
	return nil
}

func (v OSSView) Render() string {
	if v.loading {
		return "Loading..."
	}
	return v.table.Render()
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
```

- [ ] **Step 4: 提交**

```bash
git add internal/tui/views/
git commit -m "feat: add CMS, SLS, OSS views"
```

---

## Task 9: 完善主应用模型

**Files:**
- Modify: `internal/tui/app.go`

- [ ] **Step 1: 更新主应用模型**

```go
// internal/tui/app.go
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"cloud-manage/internal/tui/components"
	"cloud-manage/internal/tui/views"
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
	activeTab  Tab
	width      int
	height     int
	ready      bool
	tabs       components.Tabs
	ecsView    views.ECSView
	cmsView    views.CMSView
	slsView    views.SLSView
	ossView    views.OSSView
	search     components.Search
	status     components.Status
	keys       KeyMap
}

func NewApp() App {
	return App{
		activeTab: TabECS,
		ecsView:   views.NewECSView(),
		cmsView:   views.NewCMSView(),
		slsView:   views.NewSLSView(),
		ossView:   views.NewOSSView(),
		search:    components.NewSearch(),
		status:    components.NewStatus(),
		keys:      Keys,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
		a.updateSizes()

	case tea.KeyMsg:
		// Search mode
		if a.search.IsActive() {
			if msg.String() == "esc" {
				a.search.Deactivate()
			} else {
				cmd = a.search.Update(msg)
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}

		// Normal mode
		switch {
		case matchKey(msg, a.keys.Quit):
			return a, tea.Quit
		case matchKey(msg, a.keys.Left):
			if a.activeTab > TabECS {
				a.activeTab--
			}
		case matchKey(msg, a.keys.Right):
			if a.activeTab < TabOSS {
				a.activeTab++
			}
		case matchKey(msg, a.keys.Search):
			a.search.Activate()
		case matchKey(msg, a.keys.Refresh):
			a.refreshData()
		default:
			// Pass to active view
			switch a.activeTab {
			case TabECS:
				cmd = a.ecsView.HandleKey(msg)
			case TabCMS:
				cmd = a.cmsView.HandleKey(msg)
			case TabSLS:
				cmd = a.slsView.HandleKey(msg)
			case TabOSS:
				cmd = a.ossView.HandleKey(msg)
			}
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	// Data loaded messages
	case views.ECSLoaded, views.ECSLoadError:
		cmd = a.ecsView.Update(msg)
	case views.CMSLoaded, views.CMSLoadError:
		cmd = a.cmsView.Update(msg)
	case views.SLSLoaded, views.SLSLoadError:
		cmd = a.slsView.Update(msg)
	case views.OSSBucketsLoaded, views.OSSObjectsLoaded, views.OSSLoadError:
		cmd = a.ossView.Update(msg)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) updateSizes() {
	contentHeight := a.height - 4 // Title + Tabs + Status + Padding
	a.ecsView.SetSize(a.width, contentHeight)
	a.cmsView.SetSize(a.width, contentHeight)
	a.slsView.SetSize(a.width, contentHeight)
	a.ossView.SetSize(a.width, contentHeight)
	a.search.SetWidth(a.width)
	a.status.SetWidth(a.width)
}

func (a *App) refreshData() {
	// TODO: Implement refresh with actual credentials
	a.status.SetMessage("Refreshing...")
}

func matchKey(msg tea.KeyMsg, binding key.Binding) bool {
	for _, k := range binding.Keys() {
		if msg.String() == k {
			return true
		}
	}
	return false
}

func (a App) View() string {
	if !a.ready {
		return "Initializing..."
	}

	// Title bar
	title := TitleBarStyle.Width(a.width).Render("Cloud Manage TUI v0.0.13")

	// Tabs
	tabs := components.RenderTabs(tabNames, int(a.activeTab), a.width)

	// Content
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
	}

	// Search
	search := a.search.Render()

	// Status
	status := a.status.Render()

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", title, tabs, content, search, status)
}
```

- [ ] **Step 2: 提交**

```bash
git add internal/tui/app.go
git commit -m "feat: complete TUI application model"
```

---

## Task 10: 集成到主程序

**Files:**
- Modify: `main.go`

- [ ] **Step 1: 添加 TUI 模式检测**

在 `main.go` 的 `detectMode()` 函数中添加：

```go
// 检测 TUI 模式
if forceTUI || (len(args) > 0 && args[0] == "tui") {
    return "tui"
}
```

在 `main()` 函数的 switch 中添加：

```go
case "tui":
    runTUI()
```

添加 `runTUI()` 函数：

```go
func runTUI() {
    p := tea.NewProgram(tui.NewApp(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

- [ ] **Step 2: 更新 usage**

```go
func printUsage() {
    fmt.Println(`Usage:
  cloud-manage [flags] <service> <action> [args...]

Modes:
  --gui             Force GUI mode (requires display)
  --tui             Force TUI mode (terminal UI)
  --cli             Force CLI mode
  (auto)            Auto-detect based on display environment

Services:
  ecs               ECS 实例管理
  cms               云监控指标查询
  sls               日志服务查询
  oss               对象存储管理

Commands:
  help              显示帮助信息
  version           显示版本号
  tui               启动终端图形界面

Flags:`)
    flag.PrintDefaults()
}
```

- [ ] **Step 3: 提交**

```bash
git add main.go
git commit -m "feat: integrate TUI mode into main program"
```

---

## 验收标准

1. **功能完整**
   - ECS 实例列表/详情/操作
   - CMS 监控查看
   - SLS 日志查询
   - OSS 浏览

2. **交互体验**
   - Tab 切换流畅
   - 快捷键响应正确
   - 搜索过滤可用
   - 实时刷新工作

3. **视觉效果**
   - 颜色主题一致
   - 布局自适应终端大小
   - 状态栏显示正确

4. **构建成功**
   - `go build ./cmd/tui/` 通过
   - `go build .` 通过
