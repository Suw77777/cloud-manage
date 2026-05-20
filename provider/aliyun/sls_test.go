//go:build integration

package aliyun

import (
	"os"
	"testing"
)

func TestSLSProvider_Integration(t *testing.T) {
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	project := os.Getenv("TEST_SLS_PROJECT")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: credentials not set")
	}
	if project == "" {
		t.Skip("Skipping integration test: TEST_SLS_PROJECT not set")
	}

	provider, err := NewSLSProvider(accessKeyId, accessKeySecret, "cn-hangzhou")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	logstores, err := provider.ListLogStores(project)
	if err != nil {
		t.Fatalf("ListLogStores failed: %v", err)
	}

	t.Logf("Found %d logstores in project %s:", len(logstores), project)
	for _, ls := range logstores {
		t.Logf("  - %s", ls)
	}
}
