package aliyun

import (
	"fmt"
	"time"

	cms "github.com/alibabacloud-go/cms-20190101/v8/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// CMSProvider wraps the Aliyun CloudMonitor SDK client.
type CMSProvider struct {
	client *cms.Client
}

// NewCMSProvider creates a new CMSProvider with the given access key and region.
func NewCMSProvider(accessKeyId, accessKeySecret, region string) (*CMSProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
	}
	config.Endpoint = tea.String("https://metrics.cn-hangzhou.aliyuncs.com")

	client, err := cms.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create CMS client: %w", err)
	}
	return &CMSProvider{client: client}, nil
}

// MetricData represents a single metric data point.
type MetricData struct {
	Timestamp int64             `json:"timestamp"`
	Value     float64           `json:"value"`
	Instance  map[string]string `json:"instance"`
}

// MetricInfo represents metric metadata.
type MetricInfo struct {
	Namespace  string `json:"namespace"`
	MetricName string `json:"metricName"`
	Period     string `json:"period"`
	Statistics string `json:"statistics"`
}

// DescribeMetricLast queries the latest metric data for a specific instance.
func (p *CMSProvider) DescribeMetricLast(namespace, metricName, instanceId string, period int) ([]MetricData, error) {
	dimensions := fmt.Sprintf(`[{"instanceId":"%s"}]`, instanceId)

	request := &cms.DescribeMetricLastRequest{
		Namespace:  tea.String(namespace),
		MetricName: tea.String(metricName),
		Dimensions: tea.String(dimensions),
		Period:     tea.String(fmt.Sprintf("%d", period)),
		StartTime:  tea.String(fmt.Sprintf("%d", time.Now().Add(-1*time.Hour).UnixMilli())),
		EndTime:    tea.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
	}

	response, err := p.client.DescribeMetricLast(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeMetricLast failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	results := make([]MetricData, 0)
	if body.Datapoints != nil {
		datapoints := tea.StringValue(body.Datapoints)
		// Parse JSON array of datapoints
		if datapoints != "" && datapoints != "[]" {
			// Simple parsing - in production, use json.Unmarshal
			results = append(results, MetricData{
				Timestamp: time.Now().UnixMilli(),
				Value:     0,
				Instance:  map[string]string{"raw": datapoints},
			})
		}
	}

	return results, nil
}

// ECSMetricData holds metric data for an ECS instance.
type ECSMetricData struct {
	InstanceId   string       `json:"instanceId"`
	CPUUtilization *float64   `json:"cpuUtilization,omitempty"`
	MemoryUtilization *float64 `json:"memoryUtilization,omitempty"`
	DiskReadBPS  *float64     `json:"diskReadBps,omitempty"`
	DiskWriteBPS *float64     `json:"diskWriteBps,omitempty"`
	InternetRX   *float64     `json:"internetRx,omitempty"`
	InternetTX   *float64     `json:"internetTx,omitempty"`
	UpdateTime   string       `json:"updateTime"`
}

// GetECSMetrics retrieves key metrics for an ECS instance.
func (p *CMSProvider) GetECSMetrics(instanceId string) (*ECSMetricData, error) {
	metrics := &ECSMetricData{
		InstanceId: instanceId,
		UpdateTime: time.Now().Format(time.RFC3339),
	}

	// Query CPU utilization
	cpuVal, err := p.getMetricValue("acs_ecs_dashboard", "CPUUtilization", instanceId)
	if err == nil && cpuVal != nil {
		metrics.CPUUtilization = cpuVal
	}

	// Query Memory utilization (requires cloud monitor agent installed)
	memVal, err := p.getMetricValue("acs_ecs_dashboard", "memory_usedutilization", instanceId)
	if err == nil && memVal != nil {
		metrics.MemoryUtilization = memVal
	}

	// Query Disk read BPS
	diskRead, err := p.getMetricValue("acs_ecs_dashboard", "DiskReadBPS", instanceId)
	if err == nil && diskRead != nil {
		metrics.DiskReadBPS = diskRead
	}

	// Query Disk write BPS
	diskWrite, err := p.getMetricValue("acs_ecs_dashboard", "DiskWriteBPS", instanceId)
	if err == nil && diskWrite != nil {
		metrics.DiskWriteBPS = diskWrite
	}

	// Query Internet inbound rate
	netRX, err := p.getMetricValue("acs_ecs_dashboard", "InternetInRate", instanceId)
	if err == nil && netRX != nil {
		metrics.InternetRX = netRX
	}

	// Query Internet outbound rate
	netTX, err := p.getMetricValue("acs_ecs_dashboard", "InternetOutRate", instanceId)
	if err == nil && netTX != nil {
		metrics.InternetTX = netTX
	}

	return metrics, nil
}

// getMetricValue queries a single metric and returns the latest value.
func (p *CMSProvider) getMetricValue(namespace, metricName, instanceId string) (*float64, error) {
	dimensions := fmt.Sprintf(`[{"instanceId":"%s"}]`, instanceId)

	request := &cms.DescribeMetricLastRequest{
		Namespace:  tea.String(namespace),
		MetricName: tea.String(metricName),
		Dimensions: tea.String(dimensions),
		Period:     tea.String("60"),
		StartTime:  tea.String(fmt.Sprintf("%d", time.Now().Add(-5*time.Minute).UnixMilli())),
		EndTime:    tea.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
	}

	response, err := p.client.DescribeMetricLast(request)
	if err != nil {
		return nil, err
	}

	body := response.Body
	if body == nil || body.Datapoints == nil {
		return nil, nil
	}

	datapoints := tea.StringValue(body.Datapoints)
	if datapoints == "" || datapoints == "[]" {
		return nil, nil
	}

	// Parse the datapoints JSON to extract the value
	// Format: [{"Average":95.5,"instanceId":"xxx","timestamp":xxx}]
	var result float64
	_, err = fmt.Sscanf(datapoints, `[{"Average":%f`, &result)
	if err != nil {
		// Try Maximum
		_, err = fmt.Sscanf(datapoints, `[{"Maximum":%f`, &result)
		if err != nil {
			return nil, nil
		}
	}

	return &result, nil
}

// DescribeMetricList queries metric data with a time range for chart display.
func (p *CMSProvider) DescribeMetricList(namespace, metricName, instanceId string, startTime, endTime int64) ([]MetricData, error) {
	dimensions := fmt.Sprintf(`[{"instanceId":"%s"}]`, instanceId)

	request := &cms.DescribeMetricListRequest{
		Namespace:  tea.String(namespace),
		MetricName: tea.String(metricName),
		Dimensions: tea.String(dimensions),
		Period:     tea.String("60"),
		StartTime:  tea.String(fmt.Sprintf("%d", startTime)),
		EndTime:    tea.String(fmt.Sprintf("%d", endTime)),
	}

	response, err := p.client.DescribeMetricList(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeMetricList failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	results := make([]MetricData, 0)
	if body.Datapoints != nil {
		datapoints := tea.StringValue(body.Datapoints)
		if datapoints != "" && datapoints != "[]" {
			results = append(results, MetricData{
				Timestamp: time.Now().UnixMilli(),
				Value:     0,
				Instance:  map[string]string{"raw": datapoints},
			})
		}
	}

	return results, nil
}
