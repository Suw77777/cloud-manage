//go:build integration

package aliyun

import (
	"os"
	"testing"
)

func TestOSSProvider_Integration(t *testing.T) {
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: credentials not set")
	}

	provider, err := NewOSSProvider(accessKeyId, accessKeySecret, "cn-hangzhou")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	buckets, err := provider.ListBuckets()
	if err != nil {
		t.Fatalf("ListBuckets failed: %v", err)
	}

	t.Logf("Found %d buckets:", len(buckets))
	for _, b := range buckets {
		t.Logf("  - %s (%s)", b.Name, b.Location)
	}
}

func TestOSSProvider_ListObjects_Integration(t *testing.T) {
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	bucket := os.Getenv("TEST_OSS_BUCKET")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: credentials not set")
	}
	if bucket == "" {
		t.Skip("Skipping integration test: TEST_OSS_BUCKET not set")
	}

	provider, err := NewOSSProvider(accessKeyId, accessKeySecret, "cn-hangzhou")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	objects, _, err := provider.ListObjects(bucket, "", 10)
	if err != nil {
		t.Fatalf("ListObjects failed: %v", err)
	}

	t.Logf("Found %d objects in bucket %s:", len(objects), bucket)
	for _, obj := range objects {
		t.Logf("  - %s (%d bytes)", obj.Key, obj.Size)
	}
}
