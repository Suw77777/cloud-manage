package aliyun

import (
	"cloud-manage/provider"
	"fmt"

	vpc "github.com/alibabacloud-go/vpc-20160428/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// VPCProvider wraps the Aliyun VPC SDK client.
type VPCProvider struct {
	client *vpc.Client
}

// NewVPCProvider creates a new VPCProvider with the given access key and region.
func NewVPCProvider(accessKeyId, accessKeySecret, region string) (*VPCProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
	}
	config.Endpoint = tea.String(fmt.Sprintf("vpc.%s.aliyuncs.com", region))

	client, err := vpc.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create VPC client: %w", err)
	}
	return &VPCProvider{client: client}, nil
}

// ListVPCs lists all VPCs in the region.
func (p *VPCProvider) ListVPCs() ([]provider.VPC, error) {
	request := &vpc.DescribeVpcsRequest{
		RegionId: p.client.RegionId,
	}

	response, err := p.client.DescribeVpcs(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeVpcs failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	vpcs := make([]provider.VPC, 0)
	if body.Vpcs != nil && body.Vpcs.Vpc != nil {
		for _, v := range body.Vpcs.Vpc {
			vpcs = append(vpcs, provider.VPC{
				VpcId:        tea.StringValue(v.VpcId),
				VpcName:      tea.StringValue(v.VpcName),
				CidrBlock:    tea.StringValue(v.CidrBlock),
				Status:       tea.StringValue(v.Status),
				RegionId:     tea.StringValue(v.RegionId),
				Description:  tea.StringValue(v.Description),
				CreationTime: tea.StringValue(v.CreationTime),
			})
		}
	}

	return vpcs, nil
}

// GetVPCDetail queries detailed information for a single VPC.
func (p *VPCProvider) GetVPCDetail(vpcId string) (*provider.VPCDetail, error) {
	request := &vpc.DescribeVpcAttributeRequest{
		VpcId: tea.String(vpcId),
	}

	response, err := p.client.DescribeVpcAttribute(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeVpcAttribute failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	vswitchIds := make([]string, 0)
	if body.VSwitchIds != nil && body.VSwitchIds.VSwitchId != nil {
		for _, id := range body.VSwitchIds.VSwitchId {
			vswitchIds = append(vswitchIds, tea.StringValue(id))
		}
	}

	return &provider.VPCDetail{
		VpcId:        tea.StringValue(body.VpcId),
		VpcName:      tea.StringValue(body.VpcName),
		CidrBlock:    tea.StringValue(body.CidrBlock),
		Status:       tea.StringValue(body.Status),
		RegionId:     tea.StringValue(body.RegionId),
		Description:  tea.StringValue(body.Description),
		CreationTime: tea.StringValue(body.CreationTime),
		VSwitchIds:   vswitchIds,
	}, nil
}

// ListVSwitches lists all VSwitches in a VPC.
func (p *VPCProvider) ListVSwitches(vpcId string) ([]provider.VSwitch, error) {
	request := &vpc.DescribeVSwitchesRequest{
		VpcId: tea.String(vpcId),
	}

	response, err := p.client.DescribeVSwitches(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeVSwitches failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	vswitches := make([]provider.VSwitch, 0)
	if body.VSwitches != nil && body.VSwitches.VSwitch != nil {
		for _, vs := range body.VSwitches.VSwitch {
			vswitches = append(vswitches, provider.VSwitch{
				VSwitchId:    tea.StringValue(vs.VSwitchId),
				VSwitchName:  tea.StringValue(vs.VSwitchName),
				CidrBlock:    tea.StringValue(vs.CidrBlock),
				ZoneId:       tea.StringValue(vs.ZoneId),
				Status:       tea.StringValue(vs.Status),
				VpcId:        tea.StringValue(vs.VpcId),
				CreationTime: tea.StringValue(vs.CreationTime),
			})
		}
	}

	return vswitches, nil
}
