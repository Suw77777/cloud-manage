package aliyun

import (
	"testing"
)

func TestNewSLBProvider_EmptyCredentials(t *testing.T) {
	// Test that creating a provider with empty credentials doesn't panic
	_, err := NewSLBProvider("", "", "cn-hangzhou")
	if err != nil {
		t.Logf("Expected error for empty credentials: %v", err)
	}
}

func TestNewSLBProvider_InvalidRegion(t *testing.T) {
	// Test with invalid region
	_, err := NewSLBProvider("test-key", "test-secret", "invalid-region")
	if err != nil {
		t.Logf("Provider creation result for invalid region: %v", err)
	}
}
