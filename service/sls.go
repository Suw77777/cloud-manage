package service

import (
	"cloud-manage/provider"
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// SLSProviderFactory creates an SLSProvider given credentials and region.
type SLSProviderFactory func(accessKeyId, accessKeySecret, region string) (provider.SLSProvider, error)

// SLSService handles SLS business logic.
type SLSService struct {
	providerFactory SLSProviderFactory
}

// NewSLSService creates a new SLSService with default provider factory (cached).
func NewSLSService() *SLSService {
	return &SLSService{
		providerFactory: CachedFactory("sls", func(accessKeyId, accessKeySecret, region string) (provider.SLSProvider, error) {
			return aliyun.NewSLSProvider(accessKeyId, accessKeySecret, region)
		}),
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

// ParseTime parses time from various formats (ISO 8601, relative time, Unix timestamp).
func ParseTime(timeStr string, defaultTime time.Time) (time.Time, error) {
	if timeStr == "" {
		return defaultTime, nil
	}

	// Check if it's a relative time (e.g., "1h", "30m", "7d")
	if len(timeStr) > 1 {
		unit := timeStr[len(timeStr)-1]
		numStr := timeStr[:len(timeStr)-1]

		var duration time.Duration
		switch unit {
		case 'h':
			var hours int
			if _, err := fmt.Sscanf(numStr, "%d", &hours); err == nil {
				duration = time.Duration(hours) * time.Hour
				return time.Now().Add(-duration), nil
			}
		case 'm':
			var minutes int
			if _, err := fmt.Sscanf(numStr, "%d", &minutes); err == nil {
				duration = time.Duration(minutes) * time.Minute
				return time.Now().Add(-duration), nil
			}
		case 'd':
			var days int
			if _, err := fmt.Sscanf(numStr, "%d", &days); err == nil {
				duration = time.Duration(days) * 24 * time.Hour
				return time.Now().Add(-duration), nil
			}
		}
	}

	// Check if it's a Unix timestamp (pure digits)
	isUnixTimestamp := true
	for _, ch := range timeStr {
		if ch < '0' || ch > '9' {
			isUnixTimestamp = false
			break
		}
	}

	if isUnixTimestamp && len(timeStr) > 0 {
		var timestamp int64
		if _, err := fmt.Sscanf(timeStr, "%d", &timestamp); err == nil {
			// Check if it's in seconds or milliseconds
			if timestamp > 1e12 {
				// Milliseconds
				return time.UnixMilli(timestamp), nil
			}
			// Seconds
			return time.Unix(timestamp, 0), nil
		}
	}

	// Try ISO 8601 formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05+08:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间: %s (支持 ISO 8601, 相对时间如 1h/30m/7d, Unix 时间戳)", timeStr)
}

// ExportResult holds the result of an export operation.
type ExportResult struct {
	FilePath string `json:"filePath"`
	Count    int    `json:"count"`
	Format   string `json:"format"`
}

// ExportLogs exports logs to a file.
func (s *SLSService) ExportLogs(accessKeyId, accessKeySecret, region, project, logstore, query string, from, to time.Time, maxLines int64, format, outputPath string) (*ExportResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || project == "" || logstore == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret, region, project and logstore are required")
	}

	// Limit max lines
	if maxLines > 5000 {
		return nil, fmt.Errorf("导出数量超过限制: %d (最大 5000 条)", maxLines)
	}

	// Query logs
	result, err := s.QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query, from.Unix(), to.Unix(), maxLines)
	if err != nil {
		return nil, err
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("没有查询到日志")
	}

	// Generate output path if not specified
	if outputPath == "" {
		timestamp := time.Now().Format("20060102_150405")
		outputPath = fmt.Sprintf("sls_%s_%s_%s.%s", project, logstore, timestamp, format)
	}

	// Export based on format
	switch format {
	case "csv":
		err = exportToCSV(result.Entries, outputPath)
	case "json":
		err = exportToJSON(result.Entries, outputPath)
	default:
		return nil, fmt.Errorf("不支持的导出格式: %s (支持 csv, json)", format)
	}

	if err != nil {
		return nil, fmt.Errorf("导出失败: %w", err)
	}

	return &ExportResult{
		FilePath: outputPath,
		Count:    len(result.Entries),
		Format:   format,
	}, nil
}

// exportToCSV exports log entries to a CSV file.
func exportToCSV(entries []LogEntryAdapter, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// Add BOM for Excel compatibility
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Collect all unique keys
	keySet := make(map[string]bool)
	for _, entry := range entries {
		for key := range entry.Content {
			keySet[key] = true
		}
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}

	// Write header
	header := append([]string{"timestamp"}, keys...)
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("写入表头失败: %w", err)
	}

	// Write rows
	for _, entry := range entries {
		row := []string{time.Unix(entry.Timestamp, 0).Format("2006-01-02 15:04:05")}
		for _, key := range keys {
			row = append(row, entry.Content[key])
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("写入数据失败: %w", err)
		}
	}

	return nil
}

// exportToJSON exports log entries to a JSON file.
func exportToJSON(entries []LogEntryAdapter, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(entries); err != nil {
		return fmt.Errorf("写入 JSON 失败: %w", err)
	}

	return nil
}
