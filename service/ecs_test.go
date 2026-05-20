package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListInstances_MockProvider(t *testing.T) {
	mock := &provider.MockECSProvider{
		Instances: []provider.ECSInstance{
			{InstanceId: "i-001", InstanceName: "test-1", Status: "Running"},
			{InstanceId: "i-002", InstanceName: "test-2", Status: "Stopped"},
		},
		TotalCount: 2,
	}

	svc := NewECSServiceWithProvider(func(a, b, c string) (provider.ECSProvider, error) {
		return mock, nil
	})

	result, err := svc.ListInstances("key", "secret", "cn-hangzhou")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Instances) != 2 {
		t.Errorf("expected 2 instances, got %d", len(result.Instances))
	}
	if result.Instances[0].InstanceId != "i-001" {
		t.Errorf("expected first instance ID 'i-001', got '%s'", result.Instances[0].InstanceId)
	}
}

func TestListInstances_ProviderError(t *testing.T) {
	mock := &provider.MockECSProvider{
		Err: errors.New("API error"),
	}

	svc := NewECSServiceWithProvider(func(a, b, c string) (provider.ECSProvider, error) {
		return mock, nil
	})

	_, err := svc.ListInstances("key", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("expected error from provider")
	}
}

func TestListInstances_EmptyCredentials(t *testing.T) {
	svc := NewECSService()

	_, err := svc.ListInstances("", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("expected error for empty accessKeyId")
	}

	_, err = svc.ListInstances("key", "", "cn-hangzhou")
	if err == nil {
		t.Error("expected error for empty accessKeySecret")
	}

	_, err = svc.ListInstances("key", "secret", "")
	if err == nil {
		t.Error("expected error for empty region")
	}
}

func TestNewECSService(t *testing.T) {
	svc := NewECSService()
	if svc == nil {
		t.Error("expected non-nil ECSService")
	}
}

func TestListInstancesMultiRegion_EmptyRegions(t *testing.T) {
	svc := NewECSService()
	results := svc.ListInstancesMultiRegion("key", "secret", []string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty regions, got %d", len(results))
	}
}
