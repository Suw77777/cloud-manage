package service

import (
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
	"sync"
)

// ECSService handles ECS business logic.
type ECSService struct{}

// NewECSService creates a new ECSService.
func NewECSService() *ECSService {
	return &ECSService{}
}

// ECSInstanceAdapter is a provider-agnostic representation of an ECS instance.
// Used by app.go to avoid direct dependency on provider/aliyun.
type ECSInstanceAdapter struct {
	InstanceId   string
	InstanceName string
	Status       string
	RegionId     string
	ZoneId       string
	PublicIp     string
	PrivateIp    string
	CreationTime string
}

// ListInstancesResult holds the result of listing ECS instances.
type ListInstancesResult struct {
	Instances  []ECSInstanceAdapter `json:"instances"`
	TotalCount int32                `json:"totalCount"`
}

// RegionResult holds the query result for a single region.
type RegionResult struct {
	Region     string               `json:"region"`
	Instances  []ECSInstanceAdapter `json:"instances"`
	TotalCount int32                `json:"totalCount"`
	Error      string               `json:"error,omitempty"`
}

// toAdapters converts aliyun.ECSInstance slice to ECSInstanceAdapter slice.
func toAdapters(instances []aliyun.ECSInstance) []ECSInstanceAdapter {
	adapters := make([]ECSInstanceAdapter, 0, len(instances))
	for _, inst := range instances {
		adapters = append(adapters, ECSInstanceAdapter{
			InstanceId:   inst.InstanceId,
			InstanceName: inst.InstanceName,
			Status:       inst.Status,
			RegionId:     inst.RegionId,
			ZoneId:       inst.ZoneId,
			PublicIp:     inst.PublicIp,
			PrivateIp:    inst.PrivateIp,
			CreationTime: inst.CreationTime,
		})
	}
	return adapters
}

// ListInstances queries ECS instances for a single region.
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
		Instances:  toAdapters(instances),
		TotalCount: total,
	}, nil
}

// ListInstancesMultiRegion queries ECS instances for multiple regions concurrently.
// Each region is queried in its own goroutine. Errors for individual regions
// are captured per-region and do not fail the entire batch.
func (s *ECSService) ListInstancesMultiRegion(accessKeyId, accessKeySecret string, regions []string) []RegionResult {
	results := make([]RegionResult, len(regions))
	var wg sync.WaitGroup

	for i, region := range regions {
		wg.Add(1)
		go func(idx int, rgn string) {
			defer wg.Done()

			result, err := s.ListInstances(accessKeyId, accessKeySecret, rgn)
			if err != nil {
				results[idx] = RegionResult{
					Region: rgn,
					Error:  err.Error(),
				}
				return
			}
			results[idx] = RegionResult{
				Region:     rgn,
				Instances:  result.Instances,
				TotalCount: result.TotalCount,
			}
		}(i, region)
	}

	wg.Wait()
	return results
}
