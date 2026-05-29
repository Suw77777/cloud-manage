package aliyun

import (
	"testing"
)

func TestNewVPCProvider_EmptyCredentials(t *testing.T) {
	// Test that creating a provider with empty credentials doesn't panic
	// The actual SDK call will fail, but the provider should be created
	_, err := NewVPCProvider("", "", "cn-hangzhou")
	// We expect an error or success depending on SDK behavior
	// The important thing is it doesn't panic
	if err != nil {
		t.Logf("Expected error for empty credentials: %v", err)
	}
}

func TestNewVPCProvider_InvalidRegion(t *testing.T) {
	// Test with invalid region
	_, err := NewVPCProvider("test-key", "test-secret", "invalid-region")
	if err != nil {
		t.Logf("Provider creation result for invalid region: %v", err)
	}
}
