package provider

import (
	"errors"
	"testing"
)

func TestMockSLBProvider_ListSLBs_Success(t *testing.T) {
	mock := &MockSLBProvider{
		SLBs: []SLB{
			{LoadBalancerId: "lb-001", LoadBalancerName: "test-slb-1", Address: "1.2.3.4"},
			{LoadBalancerId: "lb-002", LoadBalancerName: "test-slb-2", Address: "5.6.7.8"},
		},
	}

	result, err := mock.ListSLBs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 SLBs, got %d", len(result))
	}
	if result[0].LoadBalancerId != "lb-001" {
		t.Errorf("expected first SLB ID 'lb-001', got '%s'", result[0].LoadBalancerId)
	}
}

func TestMockSLBProvider_ListSLBs_Error(t *testing.T) {
	mock := &MockSLBProvider{
		Err: errors.New("API error"),
	}

	_, err := mock.ListSLBs()
	if err == nil {
		t.Error("expected error")
	}
	if err.Error() != "API error" {
		t.Errorf("expected 'API error', got '%s'", err.Error())
	}
}

func TestMockSLBProvider_GetSLBDetail_Success(t *testing.T) {
	mock := &MockSLBProvider{
		SLBDetail: &SLBDetail{
			LoadBalancerId:   "lb-001",
			LoadBalancerName: "test-slb",
			Address:          "1.2.3.4",
			AddressType:      "internet",
		},
	}

	result, err := mock.GetSLBDetail("lb-001")
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

func TestMockSLBProvider_GetSLBDetail_Error(t *testing.T) {
	mock := &MockSLBProvider{
		Err: errors.New("not found"),
	}

	_, err := mock.GetSLBDetail("lb-001")
	if err == nil {
		t.Error("expected error")
	}
}

func TestMockSLBProvider_ListSLBListeners_Success(t *testing.T) {
	mock := &MockSLBProvider{
		Listeners: []SLBListener{
			{ListenerPort: 80, ListenerProtocol: "HTTP", Status: "running"},
			{ListenerPort: 443, ListenerProtocol: "HTTPS", Status: "running"},
		},
	}

	result, err := mock.ListSLBListeners("lb-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 listeners, got %d", len(result))
	}
	if result[0].ListenerPort != 80 {
		t.Errorf("expected first listener port 80, got %d", result[0].ListenerPort)
	}
}

func TestMockSLBProvider_ListSLBListeners_Error(t *testing.T) {
	mock := &MockSLBProvider{
		Err: errors.New("API error"),
	}

	_, err := mock.ListSLBListeners("lb-001")
	if err == nil {
		t.Error("expected error")
	}
}
