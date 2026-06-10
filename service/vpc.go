package service

import (
	"cloud-manage/provider"
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
)

// VPCProviderFactory 根据凭证和区域创建 VPCProvider。
type VPCProviderFactory func(accessKeyId, accessKeySecret, region string) (provider.VPCProvider, error)

// VPCService 处理 VPC 业务逻辑。
type VPCService struct {
	providerFactory VPCProviderFactory
}

// NewVPCService 创建使用默认 Provider 工厂的 VPCService（带缓存）。
func NewVPCService() *VPCService {
	return &VPCService{
		providerFactory: CachedFactory("vpc", func(accessKeyId, accessKeySecret, region string) (provider.VPCProvider, error) {
			return aliyun.NewVPCProvider(accessKeyId, accessKeySecret, region)
		}),
	}
}

// NewVPCServiceWithProvider 创建使用自定义 Provider 工厂的 VPCService（用于测试）。
func NewVPCServiceWithProvider(factory VPCProviderFactory) *VPCService {
	return &VPCService{providerFactory: factory}
}

// VPCAdapter 是与 Provider 无关的 VPC 表示。
type VPCAdapter struct {
	VpcId        string `json:"vpcId"`
	VpcName      string `json:"vpcName"`
	CidrBlock    string `json:"cidrBlock"`
	Status       string `json:"status"`
	RegionId     string `json:"regionId"`
	Description  string `json:"description"`
	CreationTime string `json:"creationTime"`
}

// VPCDetailAdapter 是与 Provider 无关的 VPC 详情表示。
type VPCDetailAdapter struct {
	VpcId          string   `json:"vpcId"`
	VpcName        string   `json:"vpcName"`
	CidrBlock      string   `json:"cidrBlock"`
	Status         string   `json:"status"`
	RegionId       string   `json:"regionId"`
	Description    string   `json:"description"`
	CreationTime   string   `json:"creationTime"`
	VSwitchIds     []string `json:"vswitchIds"`
	NatGatewayIds  []string `json:"natGatewayIds"`
	RouterTableIds []string `json:"routerTableIds"`
}

// VSwitchAdapter 是与 Provider 无关的 VSwitch 表示。
type VSwitchAdapter struct {
	VSwitchId    string `json:"vswitchId"`
	VSwitchName  string `json:"vswitchName"`
	CidrBlock    string `json:"cidrBlock"`
	ZoneId       string `json:"zoneId"`
	Status       string `json:"status"`
	VpcId        string `json:"vpcId"`
	CreationTime string `json:"creationTime"`
}

// ListVPCsResult 保存列出 VPC 的结果。
type ListVPCsResult struct {
	VPCs []VPCAdapter `json:"vpcs"`
}

// ListVSwitchesResult 保存列出 VSwitch 的结果。
type ListVSwitchesResult struct {
	VSwitches []VSwitchAdapter `json:"vswitches"`
}

// ListVPCs 列出区域内的所有 VPC。
func (s *VPCService) ListVPCs(accessKeyId, accessKeySecret, region string) (*ListVPCsResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret 和 region 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 VPC Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	vpcs, err := p.ListVPCs()
	if err != nil {
		return nil, fmt.Errorf("列出 VPC 失败: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]VPCAdapter, 0, len(vpcs))
	for _, v := range vpcs {
		adapters = append(adapters, VPCAdapter{
			VpcId:        v.VpcId,
			VpcName:      v.VpcName,
			CidrBlock:    v.CidrBlock,
			Status:       v.Status,
			RegionId:     v.RegionId,
			Description:  v.Description,
			CreationTime: v.CreationTime,
		})
	}

	return &ListVPCsResult{VPCs: adapters}, nil
}

// GetVPCDetail 查询单个 VPC 的详细信息。
func (s *VPCService) GetVPCDetail(accessKeyId, accessKeySecret, region, vpcId string) (*VPCDetailAdapter, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || vpcId == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret、region 和 vpcId 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 VPC Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	detail, err := p.GetVPCDetail(vpcId)
	if err != nil {
		return nil, fmt.Errorf("获取 VPC 详情失败: %s", security.SanitizeErrorMessage(err))
	}

	return &VPCDetailAdapter{
		VpcId:          detail.VpcId,
		VpcName:        detail.VpcName,
		CidrBlock:      detail.CidrBlock,
		Status:         detail.Status,
		RegionId:       detail.RegionId,
		Description:    detail.Description,
		CreationTime:   detail.CreationTime,
		VSwitchIds:     detail.VSwitchIds,
		NatGatewayIds:  detail.NatGatewayIds,
		RouterTableIds: detail.RouterTableIds,
	}, nil
}

// ListVSwitches 列出 VPC 中的所有虚拟交换机。
func (s *VPCService) ListVSwitches(accessKeyId, accessKeySecret, region, vpcId string) (*ListVSwitchesResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || vpcId == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret、region 和 vpcId 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 VPC Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	vswitches, err := p.ListVSwitches(vpcId)
	if err != nil {
		return nil, fmt.Errorf("列出 VSwitch 失败: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]VSwitchAdapter, 0, len(vswitches))
	for _, vs := range vswitches {
		adapters = append(adapters, VSwitchAdapter{
			VSwitchId:    vs.VSwitchId,
			VSwitchName:  vs.VSwitchName,
			CidrBlock:    vs.CidrBlock,
			ZoneId:       vs.ZoneId,
			Status:       vs.Status,
			VpcId:        vs.VpcId,
			CreationTime: vs.CreationTime,
		})
	}

	return &ListVSwitchesResult{VSwitches: adapters}, nil
}
