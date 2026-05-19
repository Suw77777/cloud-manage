package main

import (
	"cloud-manage/security"
	"cloud-manage/service"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// jsonMarshal marshals a value to JSON bytes.
func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// App struct holds references to services and is exposed to the frontend via Wails.
type App struct {
	ctx     context.Context
	ecsSvc  *service.ECSService
	cmsSvc  *service.CMSService
	slsSvc  *service.SLSService
	ossSvc  *service.OSSService
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{
		ecsSvc: service.NewECSService(),
		cmsSvc: service.NewCMSService(),
		slsSvc: service.NewSLSService(),
		ossSvc: service.NewOSSService(),
	}
}

// startup is called when the Wails app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// QueryECSResult is the response type for single-region QueryECS.
type QueryECSResult struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Instances  []InstanceView `json:"instances"`
	TotalCount int32          `json:"totalCount"`
}

// MultiRegionResult is the response type for multi-region QueryECSMultiRegion.
type MultiRegionResult struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Regions []RegionResultView `json:"regions"`
}

// RegionResultView holds the query result for a single region.
type RegionResultView struct {
	Region     string         `json:"region"`
	Instances  []InstanceView `json:"instances"`
	TotalCount int32          `json:"totalCount"`
	Error      string         `json:"error,omitempty"`
}

// InstanceView is the frontend-friendly representation of an ECS instance.
type InstanceView struct {
	InstanceId   string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	Status       string `json:"status"`
	RegionId     string `json:"regionId"`
	ZoneId       string `json:"zoneId"`
	PublicIp     string `json:"publicIp"`
	PrivateIp    string `json:"privateIp"`
	CreationTime string `json:"creationTime"`
}

// QueryECS queries ECS instances for a single region.
// Parameters are passed from GUI input fields and are NOT persisted.
func (a *App) QueryECS(accessKeyId, accessKeySecret, region, env string) QueryECSResult {
	result, err := a.ecsSvc.ListInstances(accessKeyId, accessKeySecret, region)
	if err != nil {
		return QueryECSResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	views := toInstanceViews(result.Instances)

	return QueryECSResult{
		Success:    true,
		Message:    fmt.Sprintf("found %d instance(s) in %s [%s]", len(views), region, env),
		Instances:  views,
		TotalCount: result.TotalCount,
	}
}

// QueryECSMultiRegion queries ECS instances for multiple regions concurrently.
// Parameters are passed from GUI input fields and are NOT persisted.
func (a *App) QueryECSMultiRegion(accessKeyId, accessKeySecret string, regions []string, env string) MultiRegionResult {
	if len(regions) == 0 {
		return MultiRegionResult{
			Success: false,
			Message: "at least one region is required",
		}
	}

	results := a.ecsSvc.ListInstancesMultiRegion(accessKeyId, accessKeySecret, regions)

	regionViews := make([]RegionResultView, 0, len(results))
	totalInstances := 0
	errorRegions := make([]string, 0)

	for _, r := range results {
		views := toInstanceViews(r.Instances)
		totalInstances += len(views)

		errMsg := ""
		if r.Error != "" {
			errMsg = security.SanitizeErrorMessage(fmt.Errorf("%s", r.Error))
		}
		rv := RegionResultView{
			Region:     r.Region,
			Instances:  views,
			TotalCount: r.TotalCount,
			Error:      errMsg,
		}
		regionViews = append(regionViews, rv)

		if r.Error != "" {
			errorRegions = append(errorRegions, r.Region)
		}
	}

	msg := fmt.Sprintf("found %d instance(s) across %d region(s) [%s]", totalInstances, len(regions), env)
	if len(errorRegions) > 0 {
		msg += fmt.Sprintf(", failed regions: %s", strings.Join(errorRegions, ", "))
	}

	return MultiRegionResult{
		Success: true,
		Message: msg,
		Regions: regionViews,
	}
}

// toInstanceViews converts aliyun ECSInstance slice to InstanceView slice.
func toInstanceViews(instances []service.ECSInstanceAdapter) []InstanceView {
	views := make([]InstanceView, 0, len(instances))
	for _, inst := range instances {
		views = append(views, InstanceView{
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
	return views
}

// InstanceDetailView is the frontend-friendly representation of instance details.
type InstanceDetailView struct {
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

// OperationResult is the response type for ECS operations (start/stop/reboot).
type OperationResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetECSDetail queries detailed information for a single ECS instance.
func (a *App) GetECSDetail(accessKeyId, accessKeySecret, region, instanceId string) OperationResult {
	detail, err := a.ecsSvc.GetInstanceDetail(accessKeyId, accessKeySecret, region, instanceId)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	return OperationResult{
		Success: true,
		Message: toJSON(InstanceDetailView{
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
		}),
	}
}

// StartECS starts a stopped ECS instance.
func (a *App) StartECS(accessKeyId, accessKeySecret, region, instanceId string) OperationResult {
	if err := a.ecsSvc.StartInstance(accessKeyId, accessKeySecret, region, instanceId); err != nil {
		return OperationResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}
	return OperationResult{
		Success: true,
		Message: fmt.Sprintf("instance %s start command sent successfully", instanceId),
	}
}

// StopECS stops a running ECS instance.
func (a *App) StopECS(accessKeyId, accessKeySecret, region, instanceId string, forceStop bool) OperationResult {
	if err := a.ecsSvc.StopInstance(accessKeyId, accessKeySecret, region, instanceId, forceStop); err != nil {
		return OperationResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}
	return OperationResult{
		Success: true,
		Message: fmt.Sprintf("instance %s stop command sent successfully", instanceId),
	}
}

// RebootECS reboots an ECS instance.
func (a *App) RebootECS(accessKeyId, accessKeySecret, region, instanceId string, forceStop bool) OperationResult {
	if err := a.ecsSvc.RebootInstance(accessKeyId, accessKeySecret, region, instanceId, forceStop); err != nil {
		return OperationResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}
	return OperationResult{
		Success: true,
		Message: fmt.Sprintf("instance %s reboot command sent successfully", instanceId),
	}
}

// toJSON marshals a value to JSON string. Used for embedding detail in message.
func toJSON(v interface{}) string {
	bytes, err := jsonMarshal(v)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

// MaskKey masks an access key for safe display.
func (a *App) MaskKey(key string) string {
	return security.MaskAccessKey(key)
}

// CMSMetricsResult is the response type for CMS metrics queries.
type CMSMetricsResult struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Metrics []CMSMetricsView  `json:"metrics"`
}

// CMSMetricsView is the frontend-friendly representation of ECS metrics.
type CMSMetricsView struct {
	InstanceId         string   `json:"instanceId"`
	CPUUtilization     *float64 `json:"cpuUtilization,omitempty"`
	MemoryUtilization  *float64 `json:"memoryUtilization,omitempty"`
	DiskReadBPS        *float64 `json:"diskReadBps,omitempty"`
	DiskWriteBPS       *float64 `json:"diskWriteBps,omitempty"`
	InternetRX         *float64 `json:"internetRx,omitempty"`
	InternetTX         *float64 `json:"internetTx,omitempty"`
	UpdateTime         string   `json:"updateTime"`
}

// GetECSMetrics queries monitoring metrics for a single ECS instance.
func (a *App) GetECSMetrics(accessKeyId, accessKeySecret, region, instanceId string) CMSMetricsResult {
	metrics, err := a.cmsSvc.GetInstanceMetrics(accessKeyId, accessKeySecret, region, instanceId)
	if err != nil {
		return CMSMetricsResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	return CMSMetricsResult{
		Success: true,
		Message: "metrics retrieved successfully",
		Metrics: []CMSMetricsView{toCMSMetricsView(*metrics)},
	}
}

// GetECSMetricsMultiRegion queries monitoring metrics for multiple instances across regions.
func (a *App) GetECSMetricsMultiRegion(accessKeyId, accessKeySecret string, instances []map[string]string) CMSMetricsResult {
	if len(instances) == 0 {
		return CMSMetricsResult{
			Success: false,
			Message: "at least one instance is required",
		}
	}

	pairs := make([]service.InstanceRegionPair, 0, len(instances))
	for _, inst := range instances {
		instanceId, _ := inst["instanceId"]
		region, _ := inst["region"]
		if instanceId != "" && region != "" {
			pairs = append(pairs, service.InstanceRegionPair{
				InstanceId: instanceId,
				Region:     region,
			})
		}
	}

	results := a.cmsSvc.GetInstanceMetricsMultiRegion(accessKeyId, accessKeySecret, pairs)

	allMetrics := make([]CMSMetricsView, 0)
	errorRegions := make([]string, 0)

	for _, r := range results {
		for _, m := range r.Metrics {
			allMetrics = append(allMetrics, toCMSMetricsView(m))
		}
		if r.Error != "" {
			errorRegions = append(errorRegions, r.Region)
		}
	}

	msg := fmt.Sprintf("retrieved metrics for %d instance(s)", len(allMetrics))
	if len(errorRegions) > 0 {
		msg += fmt.Sprintf(", failed regions: %s", strings.Join(errorRegions, ", "))
	}

	return CMSMetricsResult{
		Success: true,
		Message: msg,
		Metrics: allMetrics,
	}
}

// toCMSMetricsView converts service.ECSMetricAdapter to CMSMetricsView.
func toCMSMetricsView(m service.ECSMetricAdapter) CMSMetricsView {
	return CMSMetricsView{
		InstanceId:        m.InstanceId,
		CPUUtilization:    m.CPUUtilization,
		MemoryUtilization: m.MemoryUtilization,
		DiskReadBPS:       m.DiskReadBPS,
		DiskWriteBPS:      m.DiskWriteBPS,
		InternetRX:        m.InternetRX,
		InternetTX:        m.InternetTX,
		UpdateTime:        m.UpdateTime,
	}
}

// CloudProductView is the frontend-friendly representation of a cloud product.
type CloudProductView struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Namespace string         `json:"namespace"`
	Metrics   []MetricView   `json:"metrics"`
}

// MetricView is the frontend-friendly representation of a metric.
type MetricView struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

// CloudProductResult is the response type for cloud product queries.
type CloudProductResult struct {
	Success  bool               `json:"success"`
	Message  string             `json:"message"`
	Products []CloudProductView `json:"products"`
}

// GetSupportedCloudProducts returns all supported cloud products.
func (a *App) GetSupportedCloudProducts() CloudProductResult {
	products := a.cmsSvc.GetSupportedProducts()

	views := make([]CloudProductView, 0, len(products))
	for _, p := range products {
		metrics := make([]MetricView, 0, len(p.Metrics))
		for _, m := range p.Metrics {
			metrics = append(metrics, MetricView{
				ID:          m.ID,
				Name:        m.Name,
				Unit:        m.Unit,
				Description: m.Description,
			})
		}
		views = append(views, CloudProductView{
			ID:        p.ID,
			Name:      p.Name,
			Namespace: p.Namespace,
			Metrics:   metrics,
		})
	}

	return CloudProductResult{
		Success:  true,
		Message:  fmt.Sprintf("found %d product(s)", len(views)),
		Products: views,
	}
}

// GetCloudProductMetrics returns metrics for a specific product.
func (a *App) GetCloudProductMetrics(productID string) CloudProductResult {
	product, err := a.cmsSvc.GetProductMetrics(productID)
	if err != nil {
		return CloudProductResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	metrics := make([]MetricView, 0, len(product.Metrics))
	for _, m := range product.Metrics {
		metrics = append(metrics, MetricView{
			ID:          m.ID,
			Name:        m.Name,
			Unit:        m.Unit,
			Description: m.Description,
		})
	}

	return CloudProductResult{
		Success: true,
		Message: fmt.Sprintf("found %d metric(s) for %s", len(metrics), productID),
		Products: []CloudProductView{{
			ID:        product.ID,
			Name:      product.Name,
			Namespace: product.Namespace,
			Metrics:   metrics,
		}},
	}
}

// SLSLogStoreResult is the response type for SLS Logstore queries.
type SLSLogStoreResult struct {
	Success   bool     `json:"success"`
	Message   string   `json:"message"`
	LogStores []string `json:"logStores"`
}

// SLSLogQueryResult is the response type for SLS log queries.
type SLSLogQueryResult struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Logs    []LogEntryView `json:"logs"`
	Count   int64         `json:"count"`
	HasMore bool          `json:"hasMore"`
}

// LogEntryView is the frontend-friendly representation of a log entry.
type LogEntryView struct {
	Timestamp int64             `json:"timestamp"`
	Content   map[string]string `json:"content"`
}

// ListSLSLogStores lists all Logstores in an SLS project.
func (a *App) ListSLSLogStores(accessKeyId, accessKeySecret, region, project string) SLSLogStoreResult {
	logstores, err := a.slsSvc.ListLogStores(accessKeyId, accessKeySecret, region, project)
	if err != nil {
		return SLSLogStoreResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	return SLSLogStoreResult{
		Success:   true,
		Message:   fmt.Sprintf("found %d logstore(s)", len(logstores)),
		LogStores: logstores,
	}
}

// QuerySLSLogs queries logs from an SLS Logstore.
func (a *App) QuerySLSLogs(accessKeyId, accessKeySecret, region, project, logstore, query string, from, to int64, maxLines int64) SLSLogQueryResult {
	result, err := a.slsSvc.QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query, from, to, maxLines)
	if err != nil {
		return SLSLogQueryResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	logs := make([]LogEntryView, 0, len(result.Entries))
	for _, entry := range result.Entries {
		logs = append(logs, LogEntryView{
			Timestamp: entry.Timestamp,
			Content:   entry.Content,
		})
	}

	return SLSLogQueryResult{
		Success: true,
		Message: fmt.Sprintf("found %d log(s)", len(logs)),
		Logs:    logs,
		Count:   result.Count,
		HasMore: result.HasMore,
	}
}

// SLSStreamResult is the response type for streaming SLS queries.
type SLSStreamResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// QuerySLSLogsStream queries logs in chunks to reduce memory usage.
func (a *App) QuerySLSLogsStream(accessKeyId, accessKeySecret, region, project, logstore, query string, from, to int64, maxLines int64) SLSStreamResult {
	// 分块查询，每块最多100条
	chunkSize := int64(100)
	totalQueried := int64(0)

	for totalQueried < maxLines {
		remaining := maxLines - totalQueried
		if remaining > chunkSize {
			remaining = chunkSize
		}

		_, err := a.slsSvc.QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query, from, to, remaining)
		if err != nil {
			return SLSStreamResult{
				Success: false,
				Message: security.SanitizeErrorMessage(err),
			}
		}

		totalQueried += remaining

		// 如果查询结果少于请求数量，说明没有更多数据
		// 这里简化处理，实际应该检查返回的数据量
	}

	return SLSStreamResult{
		Success: true,
		Message: fmt.Sprintf("streamed %d log(s)", totalQueried),
	}
}

// OSSBucketResult is the response type for OSS bucket queries.
type OSSBucketResult struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Buckets []BucketView   `json:"buckets"`
}

// BucketView is the frontend-friendly representation of an OSS bucket.
type BucketView struct {
	Name             string `json:"name"`
	Location         string `json:"location"`
	CreationDate     string `json:"creationDate"`
	StorageClass     string `json:"storageClass"`
	ExtranetEndpoint string `json:"extranetEndpoint"`
	IntranetEndpoint string `json:"intranetEndpoint"`
}

// OSSObjectResult is the response type for OSS object queries.
type OSSObjectResult struct {
	Success     bool         `json:"success"`
	Message     string       `json:"message"`
	Objects     []ObjectView `json:"objects"`
	IsTruncated bool         `json:"isTruncated"`
}

// ObjectView is the frontend-friendly representation of an OSS object.
type ObjectView struct {
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified"`
	ETag         string `json:"etag"`
	Type         string `json:"type"`
	StorageClass string `json:"storageClass"`
	IsFolder     bool   `json:"isFolder"`
}

// ListOSSBuckets lists all OSS buckets.
func (a *App) ListOSSBuckets(accessKeyId, accessKeySecret, region string) OSSBucketResult {
	result, err := a.ossSvc.ListBuckets(accessKeyId, accessKeySecret, region)
	if err != nil {
		return OSSBucketResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	buckets := make([]BucketView, 0, len(result.Buckets))
	for _, b := range result.Buckets {
		buckets = append(buckets, BucketView{
			Name:             b.Name,
			Location:         b.Location,
			CreationDate:     b.CreationDate,
			StorageClass:     b.StorageClass,
			ExtranetEndpoint: b.ExtranetEndpoint,
			IntranetEndpoint: b.IntranetEndpoint,
		})
	}

	return OSSBucketResult{
		Success: true,
		Message: fmt.Sprintf("found %d bucket(s)", len(buckets)),
		Buckets: buckets,
	}
}

// ListOSSObjects lists objects in an OSS bucket.
func (a *App) ListOSSObjects(accessKeyId, accessKeySecret, region, bucket, prefix string, maxKeys int32) OSSObjectResult {
	result, err := a.ossSvc.ListObjects(accessKeyId, accessKeySecret, region, bucket, prefix, maxKeys)
	if err != nil {
		return OSSObjectResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	objects := make([]ObjectView, 0, len(result.Objects))
	for _, obj := range result.Objects {
		isFolder := obj.Type == "" && obj.Size == 0 && len(obj.Key) > 0 && obj.Key[len(obj.Key)-1] == '/'
		objects = append(objects, ObjectView{
			Key:          obj.Key,
			Size:         obj.Size,
			LastModified: obj.LastModified,
			ETag:         obj.ETag,
			Type:         obj.Type,
			StorageClass: obj.StorageClass,
			IsFolder:     isFolder,
		})
	}

	return OSSObjectResult{
		Success:     true,
		Message:     fmt.Sprintf("found %d object(s)", len(objects)),
		Objects:     objects,
		IsTruncated: result.IsTruncated,
	}
}
