package aliyun

import (
	"fmt"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sls "github.com/alibabacloud-go/sls-20201230/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// SLSProvider wraps the Aliyun SLS SDK client.
type SLSProvider struct {
	client *sls.Client
	region string
}

// NewSLSProvider creates a new SLSProvider with the given access key and region.
func NewSLSProvider(accessKeyId, accessKeySecret, region string) (*SLSProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
	}
	config.Endpoint = tea.String(fmt.Sprintf("%s.log.aliyuncs.com", region))

	client, err := sls.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create SLS client: %w", err)
	}
	return &SLSProvider{client: client, region: region}, nil
}

// LogStore represents an SLS Logstore.
type LogStore struct {
	LogstoreName string `json:"logstoreName"`
	TTL          int32  `json:"ttl"`
	ShardCount   int32  `json:"shardCount"`
}

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp int64             `json:"timestamp"`
	Content   map[string]string `json:"content"`
}

// ListLogStores lists all Logstores in a project.
func (p *SLSProvider) ListLogStores(project string) ([]string, error) {
	request := &sls.ListLogStoresRequest{
		Size: tea.Int32(100),
	}
	response, err := p.client.ListLogStores(tea.String(project), request)
	if err != nil {
		return nil, fmt.Errorf("ListLogStores failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return []string{}, nil
	}

	logstores := make([]string, 0)
	if body.Logstores != nil {
		for _, ls := range body.Logstores {
			logstores = append(logstores, tea.StringValue(ls))
		}
	}

	return logstores, nil
}

// GetLogs queries logs from a Logstore.
func (p *SLSProvider) GetLogs(project, logstore, query string, from, to int64, maxLines int64) ([]LogEntry, int64, error) {
	// Limit maxLines to prevent excessive memory usage
	if maxLines > 10000 {
		maxLines = 10000
	}
	if maxLines < 1 {
		maxLines = 100
	}
	request := &sls.GetLogsRequest{
		From:    tea.Int32(int32(from)),
		To:      tea.Int32(int32(to)),
		Query:   tea.String(query),
		Line:    tea.Int64(maxLines),
		Reverse: tea.Bool(false),
	}

	response, err := p.client.GetLogs(tea.String(project), tea.String(logstore), request)
	if err != nil {
		return nil, 0, fmt.Errorf("GetLogs failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return []LogEntry{}, 0, nil
	}

	entries := make([]LogEntry, 0)
	for _, log := range body {
		entry := LogEntry{
			Timestamp: time.Now().UnixMilli(),
			Content:   make(map[string]string),
		}
		for k, v := range log {
			if strVal, ok := v.(string); ok {
				entry.Content[k] = strVal
			} else {
				entry.Content[k] = fmt.Sprintf("%v", v)
			}
		}
		entries = append(entries, entry)
	}

	return entries, int64(len(entries)), nil
}

// LogHistogram represents a histogram bucket.
type LogHistogram struct {
	Count    int64  `json:"count"`
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Progress string `json:"progress"`
}

// GetHistograms gets log histogram for visualization.
func (p *SLSProvider) GetHistograms(project, logstore, query string, from, to int64) ([]LogHistogram, error) {
	request := &sls.GetHistogramsRequest{
		From:  tea.Int64(from),
		To:    tea.Int64(to),
		Query: tea.String(query),
	}

	response, err := p.client.GetHistograms(tea.String(project), tea.String(logstore), request)
	if err != nil {
		return nil, fmt.Errorf("GetHistograms failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return []LogHistogram{}, nil
	}

	histograms := make([]LogHistogram, 0)
	// Response body is a slice of histogram objects
	for _, h := range body {
		hist := LogHistogram{}
		if h.Count != nil {
			hist.Count = tea.Int64Value(h.Count)
		}
		if h.From != nil {
			hist.From = tea.Int64Value(h.From)
		}
		if h.To != nil {
			hist.To = tea.Int64Value(h.To)
		}
		if h.Progress != nil {
			hist.Progress = tea.StringValue(h.Progress)
		}
		histograms = append(histograms, hist)
	}

	return histograms, nil
}
