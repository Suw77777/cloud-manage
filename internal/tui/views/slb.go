package views

import (
	"cloud-manage/internal/handler"
	"cloud-manage/internal/tui/components"
	"cloud-manage/service"

	tea "github.com/charmbracelet/bubbletea"
)

type SLBView struct {
	table      components.Table
	detail     components.Detail
	loading    bool
	showDetail bool
	handler    *handler.SLBHandler
	width      int
	height     int
	ak         string
	sk         string
	region     string
}

func NewSLBView() SLBView {
	return SLBView{
		table:   components.NewTable([]string{"SLB ID", "Name", "Address", "Type", "Status", "VPC"}),
		detail:  components.NewDetail("SLB Detail"),
		handler: handler.NewSLBHandler(),
	}
}

func (v *SLBView) SetCredentials(ak, sk, region string) {
	v.ak = ak
	v.sk = sk
	v.region = region
}

func (v *SLBView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.table.SetSize(width, height/2)
	v.detail.SetSize(width, height/2)
}

func (v *SLBView) LoadData(accessKeyId, accessKeySecret, region string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.handler.ListSLBs(accessKeyId, accessKeySecret, region)
		if err != nil {
			return SLBLoadError{err}
		}
		return SLBLoaded{result}
	}
}

type SLBLoaded struct {
	Result *service.ListSLBsResult
}

type SLBLoadError struct {
	Err error
}

func (v *SLBView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case SLBLoaded:
		v.loading = false
		rows := make([][]string, 0)
		for _, lb := range msg.Result.SLBs {
			rows = append(rows, []string{
				lb.LoadBalancerId,
				lb.LoadBalancerName,
				lb.Address,
				lb.AddressType,
				lb.Status,
				lb.VpcId,
			})
		}
		v.table.SetRows(rows)
	case SLBLoadError:
		v.loading = false
	}
	return nil
}

func (v *SLBView) HandleKey(msg tea.KeyMsg) tea.Cmd {
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
			v.detail.SetField("Address", row[2])
			v.detail.SetField("Type", row[3])
			v.detail.SetField("Status", row[4])
			v.detail.SetField("VPC", row[5])
		}
	case "esc":
		v.showDetail = false
		v.detail.Clear()
	}
	return nil
}

func (v SLBView) Render() string {
	if v.loading {
		return "Loading..."
	}

	if v.showDetail {
		return v.detail.Render()
	}

	return v.table.Render()
}

func (v *SLBView) HandleMessage(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.KeyMsg:
		return v.HandleKey(m)
	default:
		return v.Update(msg)
	}
}
