package provider

// MockVPCProvider implements VPCProvider for testing.
// It returns pre-configured data and errors to simulate API responses.
type MockVPCProvider struct {
	VPCs      []VPC
	VPCDetail *VPCDetail
	VSwitches []VSwitch
	Err       error
}

// ListVPCs returns the mock VPC list or an error if configured.
func (m *MockVPCProvider) ListVPCs() ([]VPC, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.VPCs, nil
}

// GetVPCDetail returns the mock VPC detail or an error if configured.
func (m *MockVPCProvider) GetVPCDetail(vpcId string) (*VPCDetail, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.VPCDetail, nil
}

// ListVSwitches returns the mock VSwitch list or an error if configured.
func (m *MockVPCProvider) ListVSwitches(vpcId string) ([]VSwitch, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.VSwitches, nil
}
