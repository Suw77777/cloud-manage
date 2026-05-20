package service

import (
	"cloud-manage/provider"
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
)

// SLSProviderFactory creates an SLSProvider given credentials and region.
type SLSProviderFactory func(accessKeyId, accessKeySecret, region string) (provider.SLSProvider, error)

// SLSService handles SLS business logic.
type SLSService struct {
	providerFactory SLSProviderFactory
}

// NewSLSService creates a new SLSService with default provider factory.
func NewSLSService() *SLSService {
	return &SLSService{
		providerFactory: func(accessKeyId, accessKeySecret, region string) (provider.SLSProvider, error) {
			return aliyun.NewSLSProvider(accessKeyId, accessKeySecret, region)
		},
	}
}

// NewSLSServiceWithProvider creates a new SLSService with custom provider factory (for testing).
func NewSLSServiceWithProvider(factory SLSProviderFactory) *SLSService {
	return &SLSService{providerFactory: factory}
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

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SLS provider: %s", security.SanitizeErrorMessage(err))
	}

	logstores, err := p.ListLogStores(project)
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

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SLS provider: %s", security.SanitizeErrorMessage(err))
	}

	entries, count, err := p.GetLogs(project, logstore, query, from, to, maxLines)
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
