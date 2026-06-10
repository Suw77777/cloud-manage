package service

import (
	"cloud-manage/provider"
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
	"sync"
)

// ECSProviderFactory creates an ECSProvider given credentials and region.
type ECSProviderFactory func(accessKeyId, accessKeySecret, region string) (provider.ECSProvider, error)

// ECSService handles ECS business logic.
type ECSService struct {
	providerFactory ECSProviderFactory
}

// NewECSService creates a new ECSService with default provider factory (cached).
func NewECSService() *ECSService {
	return &ECSService{
		providerFactory: CachedFactory("ecs", func(accessKeyId, accessKeySecret, region string) (provider.ECSProvider, error) {
			return aliyun.NewECSProvider(accessKeyId, accessKeySecret, region)
		}),
	}
}

// NewECSServiceWithProvider creates a new ECSService with custom provider factory (for testing).
func NewECSServiceWithProvider(factory ECSProviderFactory) *ECSService {
	return &ECSService{providerFactory: factory}
}

// ECSInstanceAdapter is a provider-agnostic representation of an ECS instance.
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

// toAdapters converts provider.ECSInstance slice to ECSInstanceAdapter slice.
func toAdapters(instances []provider.ECSInstance) []ECSInstanceAdapter {
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
func (s *ECSService) ListInstances(accessKeyId, accessKeySecret, region string) (*ListInstancesResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret and region are required")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ECS provider: %s", security.SanitizeErrorMessage(err))
	}

	instances, total, err := p.DescribeInstances(1, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %s", security.SanitizeErrorMessage(err))
	}

	return &ListInstancesResult{
		Instances:  toAdapters(instances),
		TotalCount: total,
	}, nil
}

// InstanceDetailAdapter is a provider-agnostic representation of instance details.
type InstanceDetailAdapter struct {
	InstanceId         string   `json:"instanceId"`
	InstanceName       string   `json:"instanceName"`
	Description        string   `json:"description"`
	HostName           string   `json:"hostName"`
	Status             string   `json:"status"`
	RegionId           string   `json:"regionId"`
	ZoneId             string   `json:"zoneId"`
	InstanceType       string   `json:"instanceType"`
	Cpu                int32    `json:"cpu"`
	Memory             int32    `json:"memory"`
	ImageId            string   `json:"imageId"`
	InternetChargeType string   `json:"internetChargeType"`
	CreationTime       string   `json:"creationTime"`
	ExpiredTime        string   `json:"expiredTime"`
	StoppedMode        string   `json:"stoppedMode"`
	PublicIp           []string `json:"publicIp"`
	PrivateIp          []string `json:"privateIp"`
	SecurityGroupIds   []string `json:"securityGroupIds"`
}

// GetInstanceDetail queries detailed information for a single ECS instance.
func (s *ECSService) GetInstanceDetail(accessKeyId, accessKeySecret, region, instanceId string) (*InstanceDetailAdapter, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || instanceId == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret, region and instanceId are required")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ECS provider: %s", security.SanitizeErrorMessage(err))
	}

	detail, err := p.DescribeInstanceDetail(instanceId)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance: %s", security.SanitizeErrorMessage(err))
	}

	return &InstanceDetailAdapter{
		InstanceId:         detail.InstanceId,
		InstanceName:       detail.InstanceName,
		Description:        detail.Description,
		HostName:           detail.HostName,
		Status:             detail.Status,
		RegionId:           detail.RegionId,
		ZoneId:             detail.ZoneId,
		InstanceType:       detail.InstanceType,
		Cpu:                detail.Cpu,
		Memory:             detail.Memory,
		ImageId:            detail.ImageId,
		InternetChargeType: detail.InternetChargeType,
		CreationTime:       detail.CreationTime,
		ExpiredTime:        detail.ExpiredTime,
		StoppedMode:        detail.StoppedMode,
		PublicIp:           detail.PublicIp,
		PrivateIp:          detail.PrivateIp,
		SecurityGroupIds:   detail.SecurityGroupIds,
	}, nil
}

// StartInstance starts a stopped ECS instance.
func (s *ECSService) StartInstance(accessKeyId, accessKeySecret, region, instanceId string) error {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || instanceId == "" {
		return fmt.Errorf("accessKeyId, accessKeySecret, region and instanceId are required")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return fmt.Errorf("failed to initialize ECS provider: %s", security.SanitizeErrorMessage(err))
	}

	if err := p.StartInstance(instanceId); err != nil {
		return fmt.Errorf("failed to start instance: %s", security.SanitizeErrorMessage(err))
	}
	return nil
}

// StopInstance stops a running ECS instance.
func (s *ECSService) StopInstance(accessKeyId, accessKeySecret, region, instanceId string, forceStop bool) error {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || instanceId == "" {
		return fmt.Errorf("accessKeyId, accessKeySecret, region and instanceId are required")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return fmt.Errorf("failed to initialize ECS provider: %s", security.SanitizeErrorMessage(err))
	}

	if err := p.StopInstance(instanceId, forceStop); err != nil {
		return fmt.Errorf("failed to stop instance: %s", security.SanitizeErrorMessage(err))
	}
	return nil
}

// RebootInstance reboots an ECS instance.
func (s *ECSService) RebootInstance(accessKeyId, accessKeySecret, region, instanceId string, forceStop bool) error {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || instanceId == "" {
		return fmt.Errorf("accessKeyId, accessKeySecret, region and instanceId are required")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return fmt.Errorf("failed to initialize ECS provider: %s", security.SanitizeErrorMessage(err))
	}

	if err := p.RebootInstance(instanceId, forceStop); err != nil {
		return fmt.Errorf("failed to reboot instance: %s", security.SanitizeErrorMessage(err))
	}
	return nil
}

// ListInstancesMultiRegion queries ECS instances for multiple regions concurrently.
func (s *ECSService) ListInstancesMultiRegion(accessKeyId, accessKeySecret string, regions []string) []RegionResult {
	return s.ListInstancesMultiRegionWithConcurrency(accessKeyId, accessKeySecret, regions, 3)
}

// ListInstancesMultiRegionWithConcurrency queries ECS instances for multiple regions with concurrency control.
func (s *ECSService) ListInstancesMultiRegionWithConcurrency(accessKeyId, accessKeySecret string, regions []string, concurrency int) []RegionResult {
	if concurrency < 1 {
		concurrency = 3
	}

	results := make([]RegionResult, len(regions))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for i, region := range regions {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		go func(idx int, rgn string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

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
