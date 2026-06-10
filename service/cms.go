package service

import (
	"cloud-manage/provider"
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
	"sync"
)

// CMSProviderFactory creates a CMSProvider given credentials and region.
type CMSProviderFactory func(accessKeyId, accessKeySecret, region string) (provider.CMSProvider, error)

// CMSService handles CloudMonitor business logic.
type CMSService struct {
	providerFactory CMSProviderFactory
}

// NewCMSService creates a new CMSService with default provider factory (cached).
func NewCMSService() *CMSService {
	return &CMSService{
		providerFactory: CachedFactory("cms", func(accessKeyId, accessKeySecret, region string) (provider.CMSProvider, error) {
			return aliyun.NewCMSProvider(accessKeyId, accessKeySecret, region)
		}),
	}
}

// NewCMSServiceWithProvider creates a new CMSService with custom provider factory (for testing).
func NewCMSServiceWithProvider(factory CMSProviderFactory) *CMSService {
	return &CMSService{providerFactory: factory}
}

// CloudProduct represents a cloud product with its monitoring metrics.
type CloudProduct struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Namespace string       `json:"namespace"`
	Metrics   []MetricInfo `json:"metrics"`
}

// MetricInfo represents a monitoring metric.
type MetricInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

// GetSupportedProducts returns all supported cloud products and their default metrics.
func (s *CMSService) GetSupportedProducts() []CloudProduct {
	return []CloudProduct{
		{
			ID:        "ecs",
			Name:      "云服务器 ECS",
			Namespace: "acs_ecs_dashboard",
			Metrics: []MetricInfo{
				{ID: "CPUUtilization", Name: "CPU 使用率", Unit: "%", Description: "实例 CPU 使用率"},
				{ID: "memory_usedutilization", Name: "内存使用率", Unit: "%", Description: "实例内存使用率（需安装云监控插件）"},
				{ID: "DiskReadBPS", Name: "磁盘读速率", Unit: "B/s", Description: "磁盘读取速率"},
				{ID: "DiskWriteBPS", Name: "磁盘写速率", Unit: "B/s", Description: "磁盘写入速率"},
				{ID: "InternetInRate", Name: "公网入流量", Unit: "bps", Description: "公网网络流入速率"},
				{ID: "InternetOutRate", Name: "公网出流量", Unit: "bps", Description: "公网网络流出速率"},
				{ID: "IntranetInRate", Name: "内网入流量", Unit: "bps", Description: "内网网络流入速率"},
				{ID: "IntranetOutRate", Name: "内网出流量", Unit: "bps", Description: "内网网络流出速率"},
			},
		},
		{
			ID:        "rds",
			Name:      "云数据库 RDS",
			Namespace: "acs_rds_dashboard",
			Metrics: []MetricInfo{
				{ID: "CpuUsage", Name: "CPU 使用率", Unit: "%", Description: "数据库 CPU 使用率"},
				{ID: "MemoryUsage", Name: "内存使用率", Unit: "%", Description: "数据库内存使用率"},
				{ID: "DiskUsage", Name: "磁盘使用率", Unit: "%", Description: "磁盘空间使用率"},
				{ID: "IOPSUsage", Name: "IOPS 使用率", Unit: "%", Description: "IOPS 使用率"},
				{ID: "ConnectionUsage", Name: "连接数使用率", Unit: "%", Description: "数据库连接数使用率"},
			},
		},
		{
			ID:        "slb",
			Name:      "负载均衡 SLB",
			Namespace: "acs_slb_dashboard",
			Metrics: []MetricInfo{
				{ID: "InstanceActiveConnection", Name: "活跃连接数", Unit: "个", Description: "当前活跃连接数"},
				{ID: "InstanceNewConnection", Name: "新建连接数", Unit: "个/秒", Description: "每秒新建连接数"},
				{ID: "InstanceTrafficRX", Name: "入流量", Unit: "bytes/s", Description: "接收数据速率"},
				{ID: "InstanceTrafficTX", Name: "出流量", Unit: "bytes/s", Description: "发送数据速率"},
				{ID: "InstanceDropTraffic", Name: "丢弃流量", Unit: "bytes/s", Description: "被丢弃的数据速率"},
			},
		},
		{
			ID:        "redis",
			Name:      "云数据库 Redis",
			Namespace: "acs_kvstore",
			Metrics: []MetricInfo{
				{ID: "StandardAvgRt", Name: "平均响应时间", Unit: "ms", Description: "平均请求响应时间"},
				{ID: "StandardMaxRt", Name: "最大响应时间", Unit: "ms", Description: "最大请求响应时间"},
				{ID: "UsedMemory", Name: "已用内存", Unit: "bytes", Description: "已使用的内存"},
				{ID: "UsedMemoryRatio", Name: "内存使用率", Unit: "%", Description: "内存使用率"},
				{ID: "ConnectionUsage", Name: "连接数使用率", Unit: "%", Description: "连接数使用率"},
			},
		},
	}
}

// GetProductMetrics returns metrics for a specific product.
func (s *CMSService) GetProductMetrics(productID string) (*CloudProduct, error) {
	products := s.GetSupportedProducts()
	for _, p := range products {
		if p.ID == productID {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("unsupported product: %s, use 'cms products' to list available products", productID)
}

// ECSMetricAdapter is a provider-agnostic representation of ECS metrics.
type ECSMetricAdapter struct {
	InstanceId        string   `json:"instanceId"`
	CPUUtilization    *float64 `json:"cpuUtilization,omitempty"`
	MemoryUtilization *float64 `json:"memoryUtilization,omitempty"`
	DiskReadBPS       *float64 `json:"diskReadBps,omitempty"`
	DiskWriteBPS      *float64 `json:"diskWriteBps,omitempty"`
	InternetRX        *float64 `json:"internetRx,omitempty"`
	InternetTX        *float64 `json:"internetTx,omitempty"`
	UpdateTime        string   `json:"updateTime"`
}

// RegionMetricsResult holds the metrics result for a single region.
type RegionMetricsResult struct {
	Region  string             `json:"region"`
	Metrics []ECSMetricAdapter `json:"metrics"`
	Error   string             `json:"error,omitempty"`
}

// GetInstanceMetrics queries metrics for a single ECS instance.
func (s *CMSService) GetInstanceMetrics(accessKeyId, accessKeySecret, region, instanceId string) (*ECSMetricAdapter, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || instanceId == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret, region and instanceId are required")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CMS provider: %s", security.SanitizeErrorMessage(err))
	}

	metrics, err := p.GetECSMetrics(instanceId)
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
