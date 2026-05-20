//go:build integration

package aliyun

import (
	"os"
	"testing"
)

func getTestCredentials(t *testing.T) (string, string, string) {
	t.Helper()
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	region := os.Getenv("CLOUD_REGION")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: CLOUD_ACCESS_KEY_ID or CLOUD_ACCESS_KEY_SECRET not set")
	}
	if region == "" {
		region = "cn-hangzhou"
	}
	return accessKeyId, accessKeySecret, region
}

func TestECSProvider_Integration(t *testing.T) {
	accessKeyId, accessKeySecret, region := getTestCredentials(t)

	provider, err := NewECSProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	instances, total, err := provider.DescribeInstances(1, 10)
	if err != nil {
		t.Fatalf("DescribeInstances failed: %v", err)
	}

	t.Logf("Found %d instances (total: %d)", len(instances), total)
	for _, inst := range instances {
		t.Logf("  %s: %s (%s)", inst.InstanceId, inst.InstanceName, inst.Status)
	}
}
