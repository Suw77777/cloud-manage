package provider

// MockSLBProvider implements SLBProvider for testing.
type MockSLBProvider struct {
	SLBs      []SLB
	SLBDetail *SLBDetail
	Listeners []SLBListener
	Err       error
}

func (m *MockSLBProvider) ListSLBs() ([]SLB, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.SLBs, nil
}

func (m *MockSLBProvider) GetSLBDetail(slbId string) (*SLBDetail, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.SLBDetail, nil
}

func (m *MockSLBProvider) ListSLBListeners(slbId string) ([]SLBListener, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Listeners, nil
}
