//go:build integration

package aliyun

import (
	"os"
	"testing"
)

func TestCMSProvider_Integration(t *testing.T) {
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	instanceId := os.Getenv("TEST_INSTANCE_ID")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: credentials not set")
	}
	if instanceId == "" {
		t.Skip("Skipping integration test: TEST_INSTANCE_ID not set")
	}

	provider, err := NewCMSProvider(accessKeyId, accessKeySecret, "cn-hangzhou")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	metrics, err := provider.GetECSMetrics(instanceId)
	if err != nil {
		t.Fatalf("GetECSMetrics failed: %v", err)
	}

	t.Logf("Metrics for %s:", instanceId)
	if metrics.CPUUtilization != nil {
		t.Logf("  CPU: %.2f%%", *metrics.CPUUtilization)
	}
	if metrics.MemoryUtilization != nil {
		t.Logf("  Memory: %.2f%%", *metrics.MemoryUtilization)
	}
}
