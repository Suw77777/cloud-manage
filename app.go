package main

import (
	"cloud-manage/security"
	"cloud-manage/service"
	"context"
	"fmt"
)

// App struct holds references to services and is exposed to the frontend via Wails.
type App struct {
	ctx      context.Context
	ecsSvc   *service.ECSService
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

// QueryECSResult is the response type for QueryECS.
type QueryECSResult struct {
	Success    bool                   `json:"success"`
	Message    string                 `json:"message"`
	Instances  []InstanceView         `json:"instances"`
	TotalCount int32                  `json:"totalCount"`
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

// QueryECS is called from the frontend to query ECS instances.
// Parameters are passed from GUI input fields and are NOT persisted.
func (a *App) QueryECS(accessKeyId, accessKeySecret, region, env string) QueryECSResult {
	result, err := a.ecsSvc.ListInstances(accessKeyId, accessKeySecret, region)
	if err != nil {
		return QueryECSResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	views := make([]InstanceView, 0, len(result.Instances))
	for _, inst := range result.Instances {
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

	return QueryECSResult{
		Success:    true,
		Message:    fmt.Sprintf("found %d instance(s) in %s [%s]", len(views), region, env),
		Instances:  views,
		TotalCount: result.TotalCount,
	}
}

// MaskKey masks an access key for safe display.
func (a *App) MaskKey(key string) string {
	return security.MaskAccessKey(key)
}
