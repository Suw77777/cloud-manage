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
