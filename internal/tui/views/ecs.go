package views

import (
	"cloud-manage/internal/consts"
	"cloud-manage/internal/tui/components"
	"cloud-manage/service"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ECSView struct {
	table      components.Table
	detail     components.Detail
	loading    bool
	showDetail bool
	service    *service.ECSService
	width      int
	height     int
	ak         string
	sk         string
	region     string
}

func NewECSView() ECSView {
	return ECSView{
		table:   components.NewTable([]string{"ID", "Name", "Status", "Public IP", "Private IP"}),
		detail:  components.NewDetail("Instance Detail"),
		service: service.NewECSService(),
	}
}

func (v *ECSView) SetCredentials(ak, sk, region string) {
	v.ak = ak
	v.sk = sk
	v.region = region
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
	case "R":
		// Cycle to next region
		regions := consts.AllRegions
		for i, r := range regions {
			if r.ID == v.region {
				v.region = regions[(i+1)%len(regions)].ID
				break
			}
		}
		v.loading = true
		return v.LoadData(v.ak, v.sk, v.region)
	}
	return nil
}

func (v ECSView) Render() string {
	if v.loading {
		return fmt.Sprintf("  Loading %s...", v.region)
	}

	if v.showDetail {
		return v.detail.Render()
	}

	return fmt.Sprintf("  Region: %s (R to switch)\n\n", v.region) + v.table.Render()
}

// HandleMessage routes messages to Update or HandleKey.
func (v *ECSView) HandleMessage(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.KeyMsg:
		return v.HandleKey(m)
	default:
		return v.Update(msg)
	}
}
