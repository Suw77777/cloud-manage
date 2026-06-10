package views

import (
	"cloud-manage/internal/handler"
	"cloud-manage/internal/tui/components"
	"cloud-manage/service"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type SLSView struct {
	table   components.Table
	detail  components.Detail
	loading bool
	handler *handler.SLSHandler
	width   int
	height  int
	ak      string
	sk      string
	region  string
}

func NewSLSView() SLSView {
	return SLSView{
		table:   components.NewTable([]string{"Timestamp", "Level", "Message"}),
		detail:  components.NewDetail("Log Detail"),
		handler: handler.NewSLSHandler(),
	}
}

func (v *SLSView) SetCredentials(ak, sk, region string) {
	v.ak = ak
	v.sk = sk
	v.region = region
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
		result, err := v.handler.QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query, from, to, 100)
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
	}
	return nil
}

func (v SLSView) Render() string {
	if v.loading {
		return "Loading..."
	}
	return v.table.Render()
}

func (v *SLSView) HandleMessage(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.KeyMsg:
		return v.HandleKey(m)
	default:
		return v.Update(msg)
	}
}
