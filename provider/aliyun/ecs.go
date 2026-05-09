package aliyun

import (
	"fmt"

	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// ECSInstance represents an ECS instance for the application layer.
type ECSInstance struct {
	InstanceId   string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	Status       string `json:"status"`
	RegionId     string `json:"regionId"`
	ZoneId       string `json:"zoneId"`
	PublicIp     string `json:"publicIp"`
	PrivateIp    string `json:"privateIp"`
	CreationTime string `json:"creationTime"`
}

// ECSProvider wraps the Aliyun ECS SDK client.
type ECSProvider struct {
	client *ecs.Client
}

// NewECSProvider creates a new ECSProvider with the given access key and region.
func NewECSProvider(accessKeyId, accessKeySecret, region string) (*ECSProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
	}
	config.Endpoint = tea.String(fmt.Sprintf("ecs.%s.aliyuncs.com", region))

	client, err := ecs.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS client: %w", err)
	}
	return &ECSProvider{client: client}, nil
}

// DescribeInstances queries ECS instances with pagination.
func (p *ECSProvider) DescribeInstances(pageNumber, pageSize int32) ([]ECSInstance, int32, error) {
	request := &ecs.DescribeInstancesRequest{
		PageNumber: tea.Int32(pageNumber),
		PageSize:   tea.Int32(pageSize),
	}

	response, err := p.client.DescribeInstances(request)
	if err != nil {
		return nil, 0, fmt.Errorf("DescribeInstances failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, 0, fmt.Errorf("empty response body")
	}

	total := tea.Int32Value(body.TotalCount)

	instances := make([]ECSInstance, 0)
	if body.Instances != nil && body.Instances.Instance != nil {
		for _, inst := range body.Instances.Instance {
			publicIp := ""
			if inst.PublicIpAddress != nil && inst.PublicIpAddress.IpAddress != nil && len(inst.PublicIpAddress.IpAddress) > 0 {
				publicIp = tea.StringValue(inst.PublicIpAddress.IpAddress[0])
			}
			privateIp := ""
			if inst.VpcAttributes != nil && inst.VpcAttributes.PrivateIpAddress != nil && inst.VpcAttributes.PrivateIpAddress.IpAddress != nil && len(inst.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
				privateIp = tea.StringValue(inst.VpcAttributes.PrivateIpAddress.IpAddress[0])
			}

			instances = append(instances, ECSInstance{
				InstanceId:   tea.StringValue(inst.InstanceId),
				InstanceName: tea.StringValue(inst.InstanceName),
				Status:       tea.StringValue(inst.Status),
				RegionId:     tea.StringValue(inst.RegionId),
				ZoneId:       tea.StringValue(inst.ZoneId),
				PublicIp:     publicIp,
				PrivateIp:    privateIp,
				CreationTime: tea.StringValue(inst.CreationTime),
			})
		}
	}

	return instances, total, nil
}
