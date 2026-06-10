package handler

import (
	"cloud-manage/service"
	"fmt"
	"time"
)

// ECSListResult holds the result of listing ECS instances.
type ECSListResult struct {
	Region    string
	Instances []service.ECSInstanceAdapter
	Total     int32
	Error     string
}

// ECSHandler handles ECS operations.
type ECSHandler struct {
	svc *service.ECSService
}

// NewECSHandler creates a new ECSHandler.
func NewECSHandler() *ECSHandler {
	return &ECSHandler{svc: service.NewECSService()}
}

// ListInstances lists ECS instances for a single region.
func (h *ECSHandler) ListInstances(accessKeyId, accessKeySecret, region string) (*ECSListResult, error) {
	result, err := h.svc.ListInstances(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, err
	}
	return &ECSListResult{
		Region:    region,
		Instances: result.Instances,
		Total:     result.TotalCount,
	}, nil
}

// ListInstancesMultiRegion lists ECS instances across multiple regions.
func (h *ECSHandler) ListInstancesMultiRegion(accessKeyId, accessKeySecret string, regions []string, concurrency int) []ECSListResult {
	results := h.svc.ListInstancesMultiRegionWithConcurrency(accessKeyId, accessKeySecret, regions, concurrency)
	out := make([]ECSListResult, 0, len(results))
	for _, r := range results {
		out = append(out, ECSListResult{
			Region:    r.Region,
			Instances: r.Instances,
			Total:     r.TotalCount,
			Error:     r.Error,
		})
	}
	return out
}

// GetInstanceDetail gets detailed information about an ECS instance.
func (h *ECSHandler) GetInstanceDetail(accessKeyId, accessKeySecret, region, instanceId string) (*service.InstanceDetailAdapter, error) {
	return h.svc.GetInstanceDetail(accessKeyId, accessKeySecret, region, instanceId)
}

// StartInstance starts an ECS instance.
func (h *ECSHandler) StartInstance(accessKeyId, accessKeySecret, region, instanceId string) error {
	return h.svc.StartInstance(accessKeyId, accessKeySecret, region, instanceId)
}

// StopInstance stops an ECS instance.
func (h *ECSHandler) StopInstance(accessKeyId, accessKeySecret, region, instanceId string, force bool) error {
	return h.svc.StopInstance(accessKeyId, accessKeySecret, region, instanceId, force)
}

// RebootInstance reboots an ECS instance.
func (h *ECSHandler) RebootInstance(accessKeyId, accessKeySecret, region, instanceId string, force bool) error {
	return h.svc.RebootInstance(accessKeyId, accessKeySecret, region, instanceId, force)
}

// CMSHandler handles CloudMonitor operations.
type CMSHandler struct {
	svc *service.CMSService
}

// NewCMSHandler creates a new CMSHandler.
func NewCMSHandler() *CMSHandler {
	return &CMSHandler{svc: service.NewCMSService()}
}

// GetSupportedProducts returns supported cloud products.
func (h *CMSHandler) GetSupportedProducts() []service.CloudProduct {
	return h.svc.GetSupportedProducts()
}

// GetInstanceMetrics gets monitoring metrics for an ECS instance.
func (h *CMSHandler) GetInstanceMetrics(accessKeyId, accessKeySecret, region, instanceId string) (*service.ECSMetricAdapter, error) {
	return h.svc.GetInstanceMetrics(accessKeyId, accessKeySecret, region, instanceId)
}

// SLSHandler handles SLS operations.
type SLSHandler struct {
	svc *service.SLSService
}

// NewSLSHandler creates a new SLSHandler.
func NewSLSHandler() *SLSHandler {
	return &SLSHandler{svc: service.NewSLSService()}
}

// ListLogStores lists logstores in an SLS project.
func (h *SLSHandler) ListLogStores(accessKeyId, accessKeySecret, region, project string) ([]string, error) {
	return h.svc.ListLogStores(accessKeyId, accessKeySecret, region, project)
}

// QueryLogs queries logs from an SLS logstore.
func (h *SLSHandler) QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query string, from, to int64, maxLines int64) (*service.LogQueryResult, error) {
	return h.svc.QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query, from, to, maxLines)
}

// ExportResult holds the result of exporting SLS logs.
type ExportResult struct {
	FilePath string
	Count    int
	Format   string
}

// ExportLogs exports SLS logs to a file.
func (h *SLSHandler) ExportLogs(accessKeyId, accessKeySecret, region, project, logstore, query string, fromStr, toStr string, maxLines int64, format, outputPath string) (*ExportResult, error) {
	now := time.Now()
	from, err := service.ParseTime(fromStr, now.Add(-1*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("invalid from time: %w", err)
	}
	to, err := service.ParseTime(toStr, now)
	if err != nil {
		return nil, fmt.Errorf("invalid to time: %w", err)
	}

	result, err := h.svc.ExportLogs(accessKeyId, accessKeySecret, region, project, logstore, query, from, to, maxLines, format, outputPath)
	if err != nil {
		return nil, err
	}

	return &ExportResult{
		FilePath: result.FilePath,
		Count:    result.Count,
		Format:   result.Format,
	}, nil
}

// OSSHandler handles OSS operations.
type OSSHandler struct {
	svc *service.OSSService
}

// NewOSSHandler creates a new OSSHandler.
func NewOSSHandler() *OSSHandler {
	return &OSSHandler{svc: service.NewOSSService()}
}

// ListBuckets lists all OSS buckets.
func (h *OSSHandler) ListBuckets(accessKeyId, accessKeySecret, region string) (*service.ListBucketsResult, error) {
	return h.svc.ListBuckets(accessKeyId, accessKeySecret, region)
}

// ListObjects lists objects in an OSS bucket.
func (h *OSSHandler) ListObjects(accessKeyId, accessKeySecret, region, bucket, prefix string, maxKeys int32) (*service.ListObjectsResult, error) {
	return h.svc.ListObjects(accessKeyId, accessKeySecret, region, bucket, prefix, maxKeys)
}

// VPCHandler handles VPC operations.
type VPCHandler struct {
	svc *service.VPCService
}

// NewVPCHandler creates a new VPCHandler.
func NewVPCHandler() *VPCHandler {
	return &VPCHandler{svc: service.NewVPCService()}
}

// ListVPCs lists all VPCs.
func (h *VPCHandler) ListVPCs(accessKeyId, accessKeySecret, region string) (*service.ListVPCsResult, error) {
	return h.svc.ListVPCs(accessKeyId, accessKeySecret, region)
}

// GetVPCDetail gets detailed information about a VPC.
func (h *VPCHandler) GetVPCDetail(accessKeyId, accessKeySecret, region, vpcId string) (*service.VPCDetailAdapter, error) {
	return h.svc.GetVPCDetail(accessKeyId, accessKeySecret, region, vpcId)
}

// ListVSwitches lists VSwitches in a VPC.
func (h *VPCHandler) ListVSwitches(accessKeyId, accessKeySecret, region, vpcId string) (*service.ListVSwitchesResult, error) {
	return h.svc.ListVSwitches(accessKeyId, accessKeySecret, region, vpcId)
}

// SLBHandler handles SLB operations.
type SLBHandler struct {
	svc *service.SLBService
}

// NewSLBHandler creates a new SLBHandler.
func NewSLBHandler() *SLBHandler {
	return &SLBHandler{svc: service.NewSLBService()}
}

// ListSLBs lists all SLB instances.
func (h *SLBHandler) ListSLBs(accessKeyId, accessKeySecret, region string) (*service.ListSLBsResult, error) {
	return h.svc.ListSLBs(accessKeyId, accessKeySecret, region)
}

// GetSLBDetail gets detailed information about an SLB instance.
func (h *SLBHandler) GetSLBDetail(accessKeyId, accessKeySecret, region, slbId string) (*service.SLBDetailAdapter, error) {
	return h.svc.GetSLBDetail(accessKeyId, accessKeySecret, region, slbId)
}

// ListSLBListeners lists listeners of an SLB instance.
func (h *SLBHandler) ListSLBListeners(accessKeyId, accessKeySecret, region, slbId string) (*service.ListSLBListenersResult, error) {
	return h.svc.ListSLBListeners(accessKeyId, accessKeySecret, region, slbId)
}
