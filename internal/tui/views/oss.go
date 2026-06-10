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
	ak         string
	sk         string
	region     string
}

func NewOSSView() OSSView {
	return OSSView{
		table:   components.NewTable([]string{"Name", "Location", "Created"}),
		detail:  components.NewDetail("Bucket Detail"),
		service: service.NewOSSService(),
	}
}

func (v *OSSView) SetCredentials(ak, sk, region string) {
	v.ak = ak
	v.sk = sk
	v.region = region
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
				return v.LoadObjects(v.ak, v.sk, v.region, v.bucket)
			}
		}
	case "esc":
		if v.showBucket {
			v.showBucket = false
			return v.LoadBuckets(v.ak, v.sk, v.region)
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

func (v *OSSView) HandleMessage(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.KeyMsg:
		return v.HandleKey(m)
	default:
		return v.Update(msg)
	}
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
