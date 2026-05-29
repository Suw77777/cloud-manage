//go:build integration

package aliyun

import (
	"testing"
)

func TestSLBProvider_Integration(t *testing.T) {
	accessKeyId, accessKeySecret, region := getTestCredentials(t)

	provider, err := NewSLBProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	slbs, err := provider.ListSLBs()
	if err != nil {
		t.Fatalf("ListSLBs failed: %v", err)
	}

	t.Logf("Found %d SLBs:", len(slbs))
	for _, lb := range slbs {
		t.Logf("  %s: %s (%s)", lb.LoadBalancerId, lb.LoadBalancerName, lb.Address)
	}
}

func TestSLBProvider_GetDetail_Integration(t *testing.T) {
	accessKeyId, accessKeySecret, region := getTestCredentials(t)

	provider, err := NewSLBProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// First get an SLB ID
	slbs, err := provider.ListSLBs()
	if err != nil {
		t.Fatalf("ListSLBs failed: %v", err)
	}
	if len(slbs) == 0 {
		t.Skip("No SLBs found, skipping detail test")
	}

	slbId := slbs[0].LoadBalancerId
	detail, err := provider.GetSLBDetail(slbId)
	if err != nil {
		t.Fatalf("GetSLBDetail failed: %v", err)
	}

	t.Logf("SLB Detail: %s (%s)", detail.LoadBalancerId, detail.Address)
	t.Logf("  Address Type: %s", detail.AddressType)
}

func TestSLBProvider_ListListeners_Integration(t *testing.T) {
	accessKeyId, accessKeySecret, region := getTestCredentials(t)

	provider, err := NewSLBProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// First get an SLB ID
	slbs, err := provider.ListSLBs()
	if err != nil {
		t.Fatalf("ListSLBs failed: %v", err)
	}
	if len(slbs) == 0 {
		t.Skip("No SLBs found, skipping listeners test")
	}

	slbId := slbs[0].LoadBalancerId
	listeners, err := provider.ListSLBListeners(slbId)
	if err != nil {
		t.Fatalf("ListSLBListeners failed: %v", err)
	}

	t.Logf("Found %d listeners for SLB %s:", len(listeners), slbId)
	for _, l := range listeners {
		t.Logf("  %s:%d (%s)", l.ListenerProtocol, l.ListenerPort, l.Status)
	}
}
