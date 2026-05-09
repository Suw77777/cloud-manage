package service

import (
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
)

// ECSService handles ECS business logic.
type ECSService struct{}

// NewECSService creates a new ECSService.
func NewECSService() *ECSService {
	return &ECSService{}
}

// ListInstancesResult holds the result of listing ECS instances.
type ListInstancesResult struct {
	Instances  []aliyun.ECSInstance `json:"instances"`
	TotalCount int32                `json:"totalCount"`
}

// ListInstances queries ECS instances using the provided credentials.
// accessKeyId, accessKeySecret, and region are passed from the GUI layer.
// They are NOT stored anywhere after this call returns.
func (s *ECSService) ListInstances(accessKeyId, accessKeySecret, region string) (*ListInstancesResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret and region are required")
	}

	provider, err := aliyun.NewECSProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ECS provider: %s", security.SanitizeErrorMessage(err))
	}

	instances, total, err := provider.DescribeInstances(1, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %s", security.SanitizeErrorMessage(err))
	}

	return &ListInstancesResult{
		Instances:  instances,
		TotalCount: total,
	}, nil
}
