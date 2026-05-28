package provider

// MockSLBProvider implements SLBProvider for testing.
// It returns pre-configured data and errors to simulate API responses.
type MockSLBProvider struct {
	SLBs      []SLB
	SLBDetail *SLBDetail
	Listeners []SLBListener
	Err       error
}

// ListSLBs returns the mock SLB list or an error if configured.
func (m *MockSLBProvider) ListSLBs() ([]SLB, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.SLBs, nil
}

// GetSLBDetail returns the mock SLB detail or an error if configured.
func (m *MockSLBProvider) GetSLBDetail(slbId string) (*SLBDetail, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.SLBDetail, nil
}

// ListSLBListeners returns the mock listener list or an error if configured.
func (m *MockSLBProvider) ListSLBListeners(slbId string) ([]SLBListener, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Listeners, nil
}
