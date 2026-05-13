package service

import (
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
	"sync"
)

// CMSService handles CloudMonitor business logic.
type CMSService struct{}

// NewCMSService creates a new CMSService.
func NewCMSService() *CMSService {
	return &CMSService{}
}

// ECSMetricAdapter is a provider-agnostic representation of ECS metrics.
type ECSMetricAdapter struct {
	InstanceId         string   `json:"instanceId"`
	CPUUtilization     *float64 `json:"cpuUtilization,omitempty"`
	MemoryUtilization  *float64 `json:"memoryUtilization,omitempty"`
	DiskReadBPS        *float64 `json:"diskReadBps,omitempty"`
	DiskWriteBPS       *float64 `json:"diskWriteBps,omitempty"`
	InternetRX         *float64 `json:"internetRx,omitempty"`
	InternetTX         *float64 `json:"internetTx,omitempty"`
	UpdateTime         string   `json:"updateTime"`
}

// RegionMetricsResult holds the metrics result for a single region.
type RegionMetricsResult struct {
	Region  string              `json:"region"`
	Metrics []ECSMetricAdapter  `json:"metrics"`
	Error   string              `json:"error,omitempty"`
}

// GetInstanceMetrics queries metrics for a single ECS instance.
func (s *CMSService) GetInstanceMetrics(accessKeyId, accessKeySecret, region, instanceId string) (*ECSMetricAdapter, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || instanceId == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret, region and instanceId are required")
	}

	provider, err := aliyun.NewCMSProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CMS provider: %s", security.SanitizeErrorMessage(err))
	}

	metrics, err := provider.GetECSMetrics(instanceId)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %s", security.SanitizeErrorMessage(err))
	}

	return &ECSMetricAdapter{
		InstanceId:        metrics.InstanceId,
		CPUUtilization:    metrics.CPUUtilization,
		MemoryUtilization: metrics.MemoryUtilization,
		DiskReadBPS:       metrics.DiskReadBPS,
		DiskWriteBPS:      metrics.DiskWriteBPS,
		InternetRX:        metrics.InternetRX,
		InternetTX:        metrics.InternetTX,
		UpdateTime:        metrics.UpdateTime,
	}, nil
}

// GetInstanceMetricsMultiRegion queries metrics for multiple instances across regions.
func (s *CMSService) GetInstanceMetricsMultiRegion(accessKeyId, accessKeySecret string, instances []InstanceRegionPair) []RegionMetricsResult {
	// Group instances by region
	regionInstances := make(map[string][]string)
	for _, inst := range instances {
		regionInstances[inst.Region] = append(regionInstances[inst.Region], inst.InstanceId)
	}

	results := make([]RegionMetricsResult, 0, len(regionInstances))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for region, instanceIds := range regionInstances {
		wg.Add(1)
		go func(rgn string, ids []string) {
			defer wg.Done()

			metrics := make([]ECSMetricAdapter, 0, len(ids))
			var firstError string

			for _, id := range ids {
				result, err := s.GetInstanceMetrics(accessKeyId, accessKeySecret, rgn, id)
				if err != nil {
					if firstError == "" {
						firstError = security.SanitizeErrorMessage(err)
					}
					// Still add the instance with error
					metrics = append(metrics, ECSMetricAdapter{
						InstanceId: id,
						UpdateTime: "",
					})
					continue
				}
				metrics = append(metrics, *result)
			}

			mu.Lock()
			results = append(results, RegionMetricsResult{
				Region:  rgn,
				Metrics: metrics,
				Error:   firstError,
			})
			mu.Unlock()
		}(region, instanceIds)
	}

	wg.Wait()
	return results
}

// InstanceRegionPair holds an instance ID and its region.
type InstanceRegionPair struct {
	InstanceId string
	Region     string
}
