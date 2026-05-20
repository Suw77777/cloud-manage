package provider

// MockCMSProvider implements CMSProvider for testing.
type MockCMSProvider struct {
	Metrics *ECSMetricData
	Err     error
}

func (m *MockCMSProvider) GetECSMetrics(instanceId string) (*ECSMetricData, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Metrics, nil
}
