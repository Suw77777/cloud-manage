package service

import (
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
)

// SLSService handles SLS business logic.
type SLSService struct{}

// NewSLSService creates a new SLSService.
func NewSLSService() *SLSService {
	return &SLSService{}
}

// LogProjectAdapter is a provider-agnostic representation of an SLS project.
type LogProjectAdapter struct {
	ProjectName string `json:"projectName"`
	Description string `json:"description"`
	Region      string `json:"region"`
	CreateTime  string `json:"createTime"`
	Status      string `json:"status"`
}

// LogStoreAdapter is a provider-agnostic representation of an SLS Logstore.
type LogStoreAdapter struct {
	LogstoreName string `json:"logstoreName"`
	TTL          int32  `json:"ttl"`
	ShardCount   int32  `json:"shardCount"`
	CreateTime   string `json:"createTime"`
	ModifyTime   string `json:"modifyTime"`
}

// LogEntryAdapter is a provider-agnostic representation of a log entry.
type LogEntryAdapter struct {
	Timestamp int64             `json:"timestamp"`
	Content   map[string]string `json:"content"`
}

// LogQueryResult holds the result of a log query.
type LogQueryResult struct {
	Entries []LogEntryAdapter `json:"entries"`
	Count   int64             `json:"count"`
	HasMore bool              `json:"hasMore"`
}

// ListLogStores lists all Logstores in a project.
func (s *SLSService) ListLogStores(accessKeyId, accessKeySecret, region, project string) ([]string, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || project == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret, region and project are required")
	}

	provider, err := aliyun.NewSLSProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SLS provider: %s", security.SanitizeErrorMessage(err))
	}

	logstores, err := provider.ListLogStores(project)
	if err != nil {
		return nil, fmt.Errorf("failed to list logstores: %s", security.SanitizeErrorMessage(err))
	}

	return logstores, nil
}

// QueryLogs queries logs from a Logstore.
func (s *SLSService) QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query string, from, to int64, maxLines int64) (*LogQueryResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || project == "" || logstore == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret, region, project and logstore are required")
	}

	provider, err := aliyun.NewSLSProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SLS provider: %s", security.SanitizeErrorMessage(err))
	}

	entries, count, err := provider.GetLogs(project, logstore, query, from, to, maxLines)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]LogEntryAdapter, 0, len(entries))
	for _, entry := range entries {
		adapters = append(adapters, LogEntryAdapter{
			Timestamp: entry.Timestamp,
			Content:   entry.Content,
		})
	}

	return &LogQueryResult{
		Entries: adapters,
		Count:   count,
		HasMore: count > maxLines,
	}, nil
}
