package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListSLBs_MockProvider(t *testing.T) {
	mock := &provider.MockSLBProvider{
		SLBs: []provider.SLB{
			{LoadBalancerId: "lb-001", LoadBalancerName: "test-slb-1", Address: "1.2.3.4", Status: "active"},
			{LoadBalancerId: "lb-002", LoadBalancerName: "test-slb-2", Address: "5.6.7.8", Status: "active"},
		},
	}

	svc := NewSLBServiceWithProvider(func(a, b, c string) (provider.SLBProvider, error) {
		return mock, nil
	})

	result, err := svc.ListSLBs("key", "secret", "cn-hangzhou")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.SLBs) != 2 {
		t.Errorf("expected 2 SLBs, got %d", len(result.SLBs))
	}
	if result.SLBs[0].LoadBalancerId != "lb-001" {
		t.Errorf("expected first SLB ID 'lb-001', got '%s'", result.SLBs[0].LoadBalancerId)
	}
}

func TestListSLBs_ProviderError(t *testing.T) {
	mock := &provider.MockSLBProvider{
		Err: errors.New("API error"),
	}

	svc := NewSLBServiceWithProvider(func(a, b, c string) (provider.SLBProvider, error) {
		return mock, nil
	})

	_, err := svc.ListSLBs("key", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("expected error from provider")
	}
}

func TestListSLBs_EmptyCredentials(t *testing.T) {
	svc := NewSLBService()

	_, err := svc.ListSLBs("", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("expected error for empty accessKeyId")
	}

	_, err = svc.ListSLBs("key", "", "cn-hangzhou")
	if err == nil {
		t.Error("expected error for empty accessKeySecret")
	}

	_, err = svc.ListSLBs("key", "secret", "")
	if err == nil {
		t.Error("expected error for empty region")
	}
}

func TestNewSLBService(t *testing.T) {
	svc := NewSLBService()
	if svc == nil {
		t.Error("expected non-nil SLBService")
	}
}

func TestGetSLBDetail_MockProvider(t *testing.T) {
	mock := &provider.MockSLBProvider{
		SLBDetail: &provider.SLBDetail{
			LoadBalancerId:   "lb-001",
			LoadBalancerName: "test-slb",
			Address:          "1.2.3.4",
			AddressType:      "internet",
			Status:           "active",
			RegionId:         "cn-hangzhou",
			VpcId:            "vpc-001",
		},
	}

	svc := NewSLBServiceWithProvider(func(a, b, c string) (provider.SLBProvider, error) {
		return mock, nil
	})

	result, err := svc.GetSLBDetail("key", "secret", "cn-hangzhou", "lb-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result.LoadBalancerId != "lb-001" {
		t.Errorf("expected SLB ID 'lb-001', got '%s'", result.LoadBalancerId)
	}
	if result.Address != "1.2.3.4" {
		t.Errorf("expected address '1.2.3.4', got '%s'", result.Address)
	}
}

func TestGetSLBDetail_EmptySlbId(t *testing.T) {
	svc := NewSLBService()

	_, err := svc.GetSLBDetail("key", "secret", "cn-hangzhou", "")
	if err == nil {
		t.Error("expected error for empty slbId")
	}
}

func TestListSLBListeners_MockProvider(t *testing.T) {
	mock := &provider.MockSLBProvider{
		Listeners: []provider.SLBListener{
			{ListenerPort: 80, ListenerProtocol: "HTTP", Status: "running"},
			{ListenerPort: 443, ListenerProtocol: "HTTPS", Status: "running"},
		},
	}

	svc := NewSLBServiceWithProvider(func(a, b, c string) (provider.SLBProvider, error) {
		return mock, nil
	})

	result, err := svc.ListSLBListeners("key", "secret", "cn-hangzhou", "lb-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Listeners) != 2 {
		t.Errorf("expected 2 listeners, got %d", len(result.Listeners))
	}
	if result.Listeners[0].ListenerPort != 80 {
		t.Errorf("expected first listener port 80, got %d", result.Listeners[0].ListenerPort)
	}
}

func TestListSLBListeners_EmptySlbId(t *testing.T) {
	svc := NewSLBService()

	_, err := svc.ListSLBListeners("key", "secret", "cn-hangzhou", "")
	if err == nil {
		t.Error("expected error for empty slbId")
	}
}
