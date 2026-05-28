package service

import (
	"cloud-manage/provider"
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
)

// SLBProviderFactory 根据凭证和区域创建 SLBProvider。
type SLBProviderFactory func(accessKeyId, accessKeySecret, region string) (provider.SLBProvider, error)

// SLBService 处理 SLB 业务逻辑。
type SLBService struct {
	providerFactory SLBProviderFactory
}

// NewSLBService 创建使用默认 Provider 工厂的 SLBService。
func NewSLBService() *SLBService {
	return &SLBService{
		providerFactory: func(accessKeyId, accessKeySecret, region string) (provider.SLBProvider, error) {
			return aliyun.NewSLBProvider(accessKeyId, accessKeySecret, region)
		},
	}
}

// NewSLBServiceWithProvider 创建使用自定义 Provider 工厂的 SLBService（用于测试）。
func NewSLBServiceWithProvider(factory SLBProviderFactory) *SLBService {
	return &SLBService{providerFactory: factory}
}

// SLBAdapter 是与 Provider 无关的 SLB 表示。
type SLBAdapter struct {
	LoadBalancerId   string `json:"loadBalancerId"`
	LoadBalancerName string `json:"loadBalancerName"`
	Address          string `json:"address"`
	AddressType      string `json:"addressType"`
	Status           string `json:"status"`
	RegionId         string `json:"regionId"`
	VpcId            string `json:"vpcId"`
	CreationTime     string `json:"creationTime"`
}

// SLBDetailAdapter 是与 Provider 无关的 SLB 详情表示。
type SLBDetailAdapter struct {
	LoadBalancerId   string `json:"loadBalancerId"`
	LoadBalancerName string `json:"loadBalancerName"`
	Address          string `json:"address"`
	AddressType      string `json:"addressType"`
	Status           string `json:"status"`
	RegionId         string `json:"regionId"`
	VpcId            string `json:"vpcId"`
	VSwitchId        string `json:"vswitchId"`
	CreationTime     string `json:"creationTime"`
	ListenerCount    int    `json:"listenerCount"`
	Bandwidth        int    `json:"bandwidth"`
}

// SLBListenerAdapter 是与 Provider 无关的 SLB 监听器表示。
type SLBListenerAdapter struct {
	ListenerPort     int    `json:"listenerPort"`
	ListenerProtocol string `json:"listenerProtocol"`
	Status           string `json:"status"`
	Bandwidth        int    `json:"bandwidth"`
}

// ListSLBsResult 保存列出 SLB 的结果。
type ListSLBsResult struct {
	SLBs []SLBAdapter `json:"slbs"`
}

// ListSLBListenersResult 保存列出 SLB 监听器的结果。
type ListSLBListenersResult struct {
	Listeners []SLBListenerAdapter `json:"listeners"`
}

// ListSLBs 列出区域内的所有 SLB 实例。
func (s *SLBService) ListSLBs(accessKeyId, accessKeySecret, region string) (*ListSLBsResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret 和 region 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 SLB Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	slbs, err := p.ListSLBs()
	if err != nil {
		return nil, fmt.Errorf("列出 SLB 失败: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]SLBAdapter, 0, len(slbs))
	for _, lb := range slbs {
		adapters = append(adapters, SLBAdapter{
			LoadBalancerId:   lb.LoadBalancerId,
			LoadBalancerName: lb.LoadBalancerName,
			Address:          lb.Address,
			AddressType:      lb.AddressType,
			Status:           lb.Status,
			RegionId:         lb.RegionId,
			VpcId:            lb.VpcId,
			CreationTime:     lb.CreationTime,
		})
	}

	return &ListSLBsResult{SLBs: adapters}, nil
}

// GetSLBDetail 查询单个 SLB 实例的详细信息。
func (s *SLBService) GetSLBDetail(accessKeyId, accessKeySecret, region, slbId string) (*SLBDetailAdapter, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || slbId == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret、region 和 slbId 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 SLB Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	detail, err := p.GetSLBDetail(slbId)
	if err != nil {
		return nil, fmt.Errorf("获取 SLB 详情失败: %s", security.SanitizeErrorMessage(err))
	}

	return &SLBDetailAdapter{
		LoadBalancerId:   detail.LoadBalancerId,
		LoadBalancerName: detail.LoadBalancerName,
		Address:          detail.Address,
		AddressType:      detail.AddressType,
		Status:           detail.Status,
		RegionId:         detail.RegionId,
		VpcId:            detail.VpcId,
		VSwitchId:        detail.VSwitchId,
		CreationTime:     detail.CreationTime,
		ListenerCount:    detail.ListenerCount,
		Bandwidth:        detail.Bandwidth,
	}, nil
}

// ListSLBListeners 列出 SLB 实例上的所有监听器。
func (s *SLBService) ListSLBListeners(accessKeyId, accessKeySecret, region, slbId string) (*ListSLBListenersResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || slbId == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret、region 和 slbId 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 SLB Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	listeners, err := p.ListSLBListeners(slbId)
	if err != nil {
		return nil, fmt.Errorf("列出 SLB 监听器失败: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]SLBListenerAdapter, 0, len(listeners))
	for _, l := range listeners {
		adapters = append(adapters, SLBListenerAdapter{
			ListenerPort:     l.ListenerPort,
			ListenerProtocol: l.ListenerProtocol,
			Status:           l.Status,
			Bandwidth:        l.Bandwidth,
		})
	}

	return &ListSLBListenersResult{Listeners: adapters}, nil
}
