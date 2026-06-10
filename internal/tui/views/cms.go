package views

import (
	"cloud-manage/internal/handler"
	"cloud-manage/internal/tui/components"
	"cloud-manage/service"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type CMSView struct {
	table   components.Table
	detail  components.Detail
	loading bool
	handler *handler.CMSHandler
	width   int
	height  int
	ak      string
	sk      string
	region  string
}

func NewCMSView() CMSView {
	return CMSView{
		table:   components.NewTable([]string{"Metric", "Value", "Unit"}),
		detail:  components.NewDetail("Metric Detail"),
		handler: handler.NewCMSHandler(),
	}
}

func (v *CMSView) SetCredentials(ak, sk, region string) {
	v.ak = ak
	v.sk = sk
	v.region = region
}

func (v *CMSView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.table.SetSize(width, height/2)
	v.detail.SetSize(width, height/2)
}

func (v *CMSView) LoadData(accessKeyId, accessKeySecret, region, instanceId string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.handler.GetInstanceMetrics(accessKeyId, accessKeySecret, region, instanceId)
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

func (v *CMSView) HandleMessage(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.KeyMsg:
		return v.HandleKey(m)
	default:
		return v.Update(msg)
	}
}
