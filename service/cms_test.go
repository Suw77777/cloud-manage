package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestGetSupportedProducts(t *testing.T) {
	svc := NewCMSService()
	products := svc.GetSupportedProducts()

	if len(products) == 0 {
		t.Error("expected non-empty products list")
	}

	found := false
	for _, p := range products {
		if p.ID == "ecs" {
			found = true
			if len(p.Metrics) == 0 {
				t.Error("expected ECS to have metrics")
			}
			break
		}
	}
	if !found {
		t.Error("expected ECS product in list")
	}
}

func TestGetProductMetrics_ValidProduct(t *testing.T) {
	svc := NewCMSService()
	product, err := svc.GetProductMetrics("ecs")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if product == nil {
		t.Fatal("expected non-nil product")
	}
	if product.ID != "ecs" {
		t.Errorf("expected product ID 'ecs', got '%s'", product.ID)
	}
}

func TestGetProductMetrics_InvalidProduct(t *testing.T) {
	svc := NewCMSService()
	_, err := svc.GetProductMetrics("invalid")

	if err == nil {
		t.Error("expected error for invalid product")
	}
}

func TestGetInstanceMetrics_MockProvider(t *testing.T) {
	cpuVal := 75.5
	mock := &provider.MockCMSProvider{
		Metrics: &provider.ECSMetricData{
			InstanceId:     "i-test",
			CPUUtilization: &cpuVal,
			UpdateTime:     "2026-05-19T00:00:00Z",
		},
	}

	svc := NewCMSServiceWithProvider(func(a, b, c string) (provider.CMSProvider, error) {
		return mock, nil
	})

	result, err := svc.GetInstanceMetrics("key", "secret", "cn-hangzhou", "i-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.InstanceId != "i-test" {
		t.Errorf("expected instanceId 'i-test', got '%s'", result.InstanceId)
	}
	if result.CPUUtilization == nil || *result.CPUUtilization != 75.5 {
		t.Error("expected CPU utilization 75.5")
	}
}

func TestGetInstanceMetrics_ProviderError(t *testing.T) {
	mock := &provider.MockCMSProvider{
		Err: errors.New("API error"),
	}

	svc := NewCMSServiceWithProvider(func(a, b, c string) (provider.CMSProvider, error) {
		return mock, nil
	})

	_, err := svc.GetInstanceMetrics("key", "secret", "cn-hangzhou", "i-test")
	if err == nil {
		t.Error("expected error from provider")
	}
}

func TestGetInstanceMetrics_EmptyCredentials(t *testing.T) {
	svc := NewCMSService()
	_, err := svc.GetInstanceMetrics("", "", "", "")
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}
