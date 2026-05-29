//go:build integration

package aliyun

import (
	"testing"
)

func TestVPCProvider_Integration(t *testing.T) {
	accessKeyId, accessKeySecret, region := getTestCredentials(t)

	provider, err := NewVPCProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	vpcs, err := provider.ListVPCs()
	if err != nil {
		t.Fatalf("ListVPCs failed: %v", err)
	}

	t.Logf("Found %d VPCs:", len(vpcs))
	for _, v := range vpcs {
		t.Logf("  %s: %s (%s)", v.VpcId, v.VpcName, v.CidrBlock)
	}
}

func TestVPCProvider_GetDetail_Integration(t *testing.T) {
	accessKeyId, accessKeySecret, region := getTestCredentials(t)

	provider, err := NewVPCProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// First get a VPC ID
	vpcs, err := provider.ListVPCs()
	if err != nil {
		t.Fatalf("ListVPCs failed: %v", err)
	}
	if len(vpcs) == 0 {
		t.Skip("No VPCs found, skipping detail test")
	}

	vpcId := vpcs[0].VpcId
	detail, err := provider.GetVPCDetail(vpcId)
	if err != nil {
		t.Fatalf("GetVPCDetail failed: %v", err)
	}

	t.Logf("VPC Detail: %s (%s)", detail.VpcId, detail.CidrBlock)
	t.Logf("  VSwitches: %d", len(detail.VSwitchIds))
}

func TestVPCProvider_ListVSwitches_Integration(t *testing.T) {
	accessKeyId, accessKeySecret, region := getTestCredentials(t)

	provider, err := NewVPCProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// First get a VPC ID
	vpcs, err := provider.ListVPCs()
	if err != nil {
		t.Fatalf("ListVPCs failed: %v", err)
	}
	if len(vpcs) == 0 {
		t.Skip("No VPCs found, skipping VSwitch test")
	}

	vpcId := vpcs[0].VpcId
	vswitches, err := provider.ListVSwitches(vpcId)
	if err != nil {
		t.Fatalf("ListVSwitches failed: %v", err)
	}

	t.Logf("Found %d VSwitches in VPC %s:", len(vswitches), vpcId)
	for _, vs := range vswitches {
		t.Logf("  %s: %s (%s)", vs.VSwitchId, vs.VSwitchName, vs.CidrBlock)
	}
}
