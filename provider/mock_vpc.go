package provider

// MockVPCProvider implements VPCProvider for testing.
type MockVPCProvider struct {
	VPCs      []VPC
	VPCDetail *VPCDetail
	VSwitches []VSwitch
	Err       error
}

func (m *MockVPCProvider) ListVPCs() ([]VPC, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.VPCs, nil
}

func (m *MockVPCProvider) GetVPCDetail(vpcId string) (*VPCDetail, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.VPCDetail, nil
}

func (m *MockVPCProvider) ListVSwitches(vpcId string) ([]VSwitch, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.VSwitches, nil
}
