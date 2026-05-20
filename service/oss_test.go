package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListBuckets_MockProvider(t *testing.T) {
	mock := &provider.MockOSSProvider{
		Buckets: []provider.OSSBucket{
			{Name: "bucket1", Location: "oss-cn-hangzhou"},
			{Name: "bucket2", Location: "oss-cn-shenzhen"},
		},
	}

	svc := NewOSSServiceWithProvider(func(a, b, c string) (provider.OSSProvider, error) {
		return mock, nil
	})

	result, err := svc.ListBuckets("key", "secret", "cn-hangzhou")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Buckets) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(result.Buckets))
	}
}

func TestListBuckets_EmptyCredentials(t *testing.T) {
	svc := NewOSSService()
	_, err := svc.ListBuckets("", "", "")
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}

func TestListObjects_MockProvider(t *testing.T) {
	mock := &provider.MockOSSProvider{
		Objects: []provider.OSSObject{
			{Key: "file1.txt", Size: 100},
			{Key: "file2.txt", Size: 200},
		},
	}

	svc := NewOSSServiceWithProvider(func(a, b, c string) (provider.OSSProvider, error) {
		return mock, nil
	})

	result, err := svc.ListObjects("key", "secret", "cn-hangzhou", "test-bucket", "", 100)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Objects) != 2 {
		t.Errorf("expected 2 objects, got %d", len(result.Objects))
	}
}

func TestListObjects_EmptyCredentials(t *testing.T) {
	svc := NewOSSService()
	_, err := svc.ListObjects("", "", "", "", "", 0)
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}

func TestDetectBucketRegion(t *testing.T) {
	mock := &provider.MockOSSProvider{
		Buckets: []provider.OSSBucket{
			{Name: "test-bucket", Location: "oss-cn-shenzhen"},
		},
	}

	svc := NewOSSServiceWithProvider(func(a, b, c string) (provider.OSSProvider, error) {
		return mock, nil
	})

	region, err := svc.DetectBucketRegion("key", "secret", "test-bucket")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if region != "cn-shenzhen" {
		t.Errorf("expected region 'cn-shenzhen', got '%s'", region)
	}
}

func TestDetectBucketRegion_NotFound(t *testing.T) {
	mock := &provider.MockOSSProvider{
		Buckets: []provider.OSSBucket{},
	}

	svc := NewOSSServiceWithProvider(func(a, b, c string) (provider.OSSProvider, error) {
		return mock, nil
	})

	_, err := svc.DetectBucketRegion("key", "secret", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent bucket")
	}
}

func TestListBuckets_ProviderError(t *testing.T) {
	mock := &provider.MockOSSProvider{
		Err: errors.New("API error"),
	}

	svc := NewOSSServiceWithProvider(func(a, b, c string) (provider.OSSProvider, error) {
		return mock, nil
	})

	_, err := svc.ListBuckets("key", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("expected error from provider")
	}
}
