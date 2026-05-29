package aliyun

import (
	"cloud-manage/provider"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	slb "github.com/alibabacloud-go/slb-20140515/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

// SLBProvider wraps the Aliyun SLB SDK client.
type SLBProvider struct {
	client *slb.Client
}

// NewSLBProvider creates a new SLBProvider with the given access key and region.
func NewSLBProvider(accessKeyId, accessKeySecret, region string) (*SLBProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
	}
	config.Endpoint = tea.String(fmt.Sprintf("slb.%s.aliyuncs.com", region))

	client, err := slb.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create SLB client: %w", err)
	}
	return &SLBProvider{client: client}, nil
}

// ListSLBs lists all SLB instances in the region.
func (p *SLBProvider) ListSLBs() ([]provider.SLB, error) {
	region := tea.StringValue(p.client.RegionId)
	request := &slb.DescribeLoadBalancersRequest{
		RegionId: tea.String(region),
	}

	response, err := p.client.DescribeLoadBalancers(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeLoadBalancers failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	slbs := make([]provider.SLB, 0)
	if body.LoadBalancers != nil && body.LoadBalancers.LoadBalancer != nil {
		for _, lb := range body.LoadBalancers.LoadBalancer {
			slbs = append(slbs, provider.SLB{
				LoadBalancerId:   tea.StringValue(lb.LoadBalancerId),
				LoadBalancerName: tea.StringValue(lb.LoadBalancerName),
				Address:          tea.StringValue(lb.Address),
				AddressType:      tea.StringValue(lb.AddressType),
				RegionId:         tea.StringValue(lb.RegionId),
				VpcId:            tea.StringValue(lb.VpcId),
				CreationTime:     tea.StringValue(lb.CreateTime),
			})
		}
	}

	return slbs, nil
}

// GetSLBDetail queries detailed information for a single SLB instance.
func (p *SLBProvider) GetSLBDetail(slbId string) (*provider.SLBDetail, error) {
	// Use list API to get details since DescribeLoadBalancerAttribute has different request format
	slbs, err := p.ListSLBs()
	if err != nil {
		return nil, err
	}

	for _, lb := range slbs {
		if lb.LoadBalancerId == slbId {
			return &provider.SLBDetail{
				LoadBalancerId:   lb.LoadBalancerId,
				LoadBalancerName: lb.LoadBalancerName,
				Address:          lb.Address,
				AddressType:      lb.AddressType,
				RegionId:         lb.RegionId,
				VpcId:            lb.VpcId,
				CreationTime:     lb.CreationTime,
			}, nil
		}
	}

	return nil, fmt.Errorf("SLB %s not found", slbId)
}

// ListSLBListeners lists all listeners on an SLB instance.
func (p *SLBProvider) ListSLBListeners(slbId string) ([]provider.SLBListener, error) {
	request := &slb.DescribeLoadBalancerListenersRequest{
		LoadBalancerId: []*string{tea.String(slbId)},
	}

	response, err := p.client.DescribeLoadBalancerListeners(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeLoadBalancerListeners failed: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	listeners := make([]provider.SLBListener, 0)
	if body.Listeners != nil {
		for _, l := range body.Listeners {
			listeners = append(listeners, provider.SLBListener{
				ListenerPort:     int(tea.Int32Value(l.ListenerPort)),
				ListenerProtocol: tea.StringValue(l.ListenerProtocol),
				Status:           tea.StringValue(l.Status),
				Bandwidth:        int(tea.Int32Value(l.Bandwidth)),
			})
		}
	}

	return listeners, nil
}
