package main

import (
	"cloud-manage/security"
	"cloud-manage/service"
	"context"
	"fmt"
	"strings"
)

// App struct holds references to services and is exposed to the frontend via Wails.
type App struct {
	ctx    context.Context
	ecsSvc *service.ECSService
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{
		ecsSvc: service.NewECSService(),
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

// MaskKey masks an access key for safe display.
func (a *App) MaskKey(key string) string {
	return security.MaskAccessKey(key)
}
