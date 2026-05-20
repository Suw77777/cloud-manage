package provider

// MockECSProvider implements ECSProvider for testing.
type MockECSProvider struct {
	Instances      []ECSInstance
	TotalCount     int32
	InstanceDetail *InstanceDetail
	Err            error
	StartErr       error
	StopErr        error
	RebootErr      error
}

func (m *MockECSProvider) DescribeInstances(pageNumber, pageSize int32) ([]ECSInstance, int32, error) {
	if m.Err != nil {
		return nil, 0, m.Err
	}
	return m.Instances, m.TotalCount, nil
}

func (m *MockECSProvider) DescribeInstanceDetail(instanceId string) (*InstanceDetail, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.InstanceDetail, nil
}

func (m *MockECSProvider) StartInstance(instanceId string) error {
	return m.StartErr
}

func (m *MockECSProvider) StopInstance(instanceId string, forceStop bool) error {
	return m.StopErr
}

func (m *MockECSProvider) RebootInstance(instanceId string, forceStop bool) error {
	return m.RebootErr
}
