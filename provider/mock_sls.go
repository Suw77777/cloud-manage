package provider

// MockSLSProvider implements SLSProvider for testing.
type MockSLSProvider struct {
	LogStores []string
	Logs      []LogEntry
	Count     int64
	Err       error
}

func (m *MockSLSProvider) ListLogStores(project string) ([]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.LogStores, nil
}

func (m *MockSLSProvider) GetLogs(project, logstore, query string, from, to int64, maxLines int64) ([]LogEntry, int64, error) {
	if m.Err != nil {
		return nil, 0, m.Err
	}
	return m.Logs, m.Count, nil
}
