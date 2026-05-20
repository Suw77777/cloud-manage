package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListLogStores_MockProvider(t *testing.T) {
	mock := &provider.MockSLSProvider{
		LogStores: []string{"logstore1", "logstore2"},
	}

	svc := NewSLSServiceWithProvider(func(a, b, c string) (provider.SLSProvider, error) {
		return mock, nil
	})

	result, err := svc.ListLogStores("key", "secret", "cn-hangzhou", "test-project")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 logstores, got %d", len(result))
	}
}

func TestListLogStores_EmptyCredentials(t *testing.T) {
	svc := NewSLSService()
	_, err := svc.ListLogStores("", "", "", "")
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}

func TestQueryLogs_MockProvider(t *testing.T) {
	mock := &provider.MockSLSProvider{
		Logs: []provider.LogEntry{
			{Timestamp: 1000, Content: map[string]string{"level": "INFO"}},
			{Timestamp: 2000, Content: map[string]string{"level": "ERROR"}},
		},
		Count: 2,
	}

	svc := NewSLSServiceWithProvider(func(a, b, c string) (provider.SLSProvider, error) {
		return mock, nil
	})

	result, err := svc.QueryLogs("key", "secret", "cn-hangzhou", "project", "logstore", "", 0, 100, 100)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}
	if len(result.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Entries))
	}
}

func TestQueryLogs_EmptyCredentials(t *testing.T) {
	svc := NewSLSService()
	_, err := svc.QueryLogs("", "", "", "", "", "", 0, 0, 0)
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}

func TestQueryLogs_ProviderError(t *testing.T) {
	mock := &provider.MockSLSProvider{
		Err: errors.New("API error"),
	}

	svc := NewSLSServiceWithProvider(func(a, b, c string) (provider.SLSProvider, error) {
		return mock, nil
	})

	_, err := svc.QueryLogs("key", "secret", "cn-hangzhou", "project", "logstore", "", 0, 100, 100)
	if err == nil {
		t.Error("expected error from provider")
	}
}
