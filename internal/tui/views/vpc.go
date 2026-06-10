package views

import (
	"cloud-manage/internal/handler"
	"cloud-manage/internal/tui/components"
	"cloud-manage/service"

	tea "github.com/charmbracelet/bubbletea"
)

type VPCView struct {
	table      components.Table
	detail     components.Detail
	loading    bool
	showDetail bool
	handler    *handler.VPCHandler
	width      int
	height     int
	ak         string
	sk         string
	region     string
}

func NewVPCView() VPCView {
	return VPCView{
		table:   components.NewTable([]string{"VPC ID", "Name", "CIDR", "Status", "Region"}),
		detail:  components.NewDetail("VPC Detail"),
		handler: handler.NewVPCHandler(),
	}
}

func (v *VPCView) SetCredentials(ak, sk, region string) {
	v.ak = ak
	v.sk = sk
	v.region = region
}

func (v *VPCView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.table.SetSize(width, height/2)
	v.detail.SetSize(width, height/2)
}

func (v *VPCView) LoadData(accessKeyId, accessKeySecret, region string) tea.Cmd {
	return func() tea.Msg {
		result, err := v.handler.ListVPCs(accessKeyId, accessKeySecret, region)
		if err != nil {
			return VPCLoadError{err}
		}
		return VPCLoaded{result}
	}
}

type VPCLoaded struct {
	Result *service.ListVPCsResult
}

type VPCLoadError struct {
	Err error
}

func (v *VPCView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case VPCLoaded:
		v.loading = false
		rows := make([][]string, 0)
		for _, vpc := range msg.Result.VPCs {
			rows = append(rows, []string{
				vpc.VpcId,
				vpc.VpcName,
				vpc.CidrBlock,
				vpc.Status,
				vpc.RegionId,
			})
		}
		v.table.SetRows(rows)
	case VPCLoadError:
		v.loading = false
	}
	return nil
}

func (v *VPCView) HandleKey(msg tea.KeyMsg) tea.Cmd {
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
			v.detail.SetField("CIDR", row[2])
			v.detail.SetField("Status", row[3])
			v.detail.SetField("Region", row[4])
		}
	case "esc":
		v.showDetail = false
		v.detail.Clear()
	}
	return nil
}

func (v VPCView) Render() string {
	if v.loading {
		return "Loading..."
	}

	if v.showDetail {
		return v.detail.Render()
	}

	return v.table.Render()
}

func (v *VPCView) HandleMessage(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.KeyMsg:
		return v.HandleKey(m)
	default:
		return v.Update(msg)
	}
}
