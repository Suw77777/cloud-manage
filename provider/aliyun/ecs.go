package aliyun

import (
	"cloud-manage/provider"
	"fmt"

	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

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

// DescribeInstanceDetail queries detailed information for a single ECS instance.
func (p *ECSProvider) DescribeInstanceDetail(instanceId string) (*provider.InstanceDetail, error) {
	request := &ecs.DescribeInstanceAttributeRequest{
		InstanceId: tea.String(instanceId),
	}

	response, err := p.client.DescribeInstanceAttribute(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeInstanceAttribute failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	publicIps := make([]string, 0)
	if body.PublicIpAddress != nil && body.PublicIpAddress.IpAddress != nil {
		for _, ip := range body.PublicIpAddress.IpAddress {
			publicIps = append(publicIps, tea.StringValue(ip))
		}
	}

	privateIps := make([]string, 0)
	if body.VpcAttributes != nil && body.VpcAttributes.PrivateIpAddress != nil && body.VpcAttributes.PrivateIpAddress.IpAddress != nil {
		for _, ip := range body.VpcAttributes.PrivateIpAddress.IpAddress {
			privateIps = append(privateIps, tea.StringValue(ip))
		}
	}

	securityGroups := make([]string, 0)
	if body.SecurityGroupIds != nil && body.SecurityGroupIds.SecurityGroupId != nil {
		for _, sg := range body.SecurityGroupIds.SecurityGroupId {
			securityGroups = append(securityGroups, tea.StringValue(sg))
		}
	}

	return &provider.InstanceDetail{
		InstanceId:         tea.StringValue(body.InstanceId),
		InstanceName:       tea.StringValue(body.InstanceName),
		Description:        tea.StringValue(body.Description),
		HostName:           tea.StringValue(body.HostName),
		Status:             tea.StringValue(body.Status),
		RegionId:           tea.StringValue(body.RegionId),
		ZoneId:             tea.StringValue(body.ZoneId),
		InstanceType:       tea.StringValue(body.InstanceType),
		Cpu:                tea.Int32Value(body.Cpu),
		Memory:             tea.Int32Value(body.Memory),
		ImageId:            tea.StringValue(body.ImageId),
		InternetChargeType: tea.StringValue(body.InternetChargeType),
		CreationTime:       tea.StringValue(body.CreationTime),
		ExpiredTime:        tea.StringValue(body.ExpiredTime),
		StoppedMode:        tea.StringValue(body.StoppedMode),
		PublicIp:           publicIps,
		PrivateIp:          privateIps,
		SecurityGroupIds:   securityGroups,
	}, nil
}

// StartInstance starts a stopped ECS instance.
func (p *ECSProvider) StartInstance(instanceId string) error {
	request := &ecs.StartInstanceRequest{
		InstanceId: tea.String(instanceId),
	}
	_, err := p.client.StartInstance(request)
	if err != nil {
		return fmt.Errorf("StartInstance failed: %w", err)
	}
	return nil
}

// StopInstance stops a running ECS instance.
func (p *ECSProvider) StopInstance(instanceId string, forceStop bool) error {
	request := &ecs.StopInstanceRequest{
		InstanceId:  tea.String(instanceId),
		ConfirmStop: tea.Bool(forceStop),
		ForceStop:   tea.Bool(forceStop),
	}
	_, err := p.client.StopInstance(request)
	if err != nil {
		return fmt.Errorf("StopInstance failed: %w", err)
	}
	return nil
}

// RebootInstance reboots an ECS instance.
func (p *ECSProvider) RebootInstance(instanceId string, forceStop bool) error {
	request := &ecs.RebootInstanceRequest{
		InstanceId: tea.String(instanceId),
		ForceStop:  tea.Bool(forceStop),
	}
	_, err := p.client.RebootInstance(request)
	if err != nil {
		return fmt.Errorf("RebootInstance failed: %w", err)
	}
	return nil
}

// DescribeInstances queries ECS instances with pagination.
func (p *ECSProvider) DescribeInstances(pageNumber, pageSize int32) ([]provider.ECSInstance, int32, error) {
	region := tea.StringValue(p.client.RegionId)
	request := &ecs.DescribeInstancesRequest{
		RegionId:   tea.String(region),
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

	instances := make([]provider.ECSInstance, 0)
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

			instances = append(instances, provider.ECSInstance{
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
