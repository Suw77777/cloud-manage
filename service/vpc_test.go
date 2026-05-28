package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListVPCs_MockProvider(t *testing.T) {
	mock := &provider.MockVPCProvider{
		VPCs: []provider.VPC{
			{VpcId: "vpc-001", VpcName: "test-vpc-1", CidrBlock: "10.0.0.0/16", Status: "Available"},
			{VpcId: "vpc-002", VpcName: "test-vpc-2", CidrBlock: "172.16.0.0/12", Status: "Available"},
		},
	}

	svc := NewVPCServiceWithProvider(func(a, b, c string) (provider.VPCProvider, error) {
		return mock, nil
	})

	result, err := svc.ListVPCs("key", "secret", "cn-hangzhou")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.VPCs) != 2 {
		t.Errorf("expected 2 VPCs, got %d", len(result.VPCs))
	}
	if result.VPCs[0].VpcId != "vpc-001" {
		t.Errorf("expected first VPC ID 'vpc-001', got '%s'", result.VPCs[0].VpcId)
	}
}

func TestListVPCs_ProviderError(t *testing.T) {
	mock := &provider.MockVPCProvider{
		Err: errors.New("API error"),
	}

	svc := NewVPCServiceWithProvider(func(a, b, c string) (provider.VPCProvider, error) {
		return mock, nil
	})

	_, err := svc.ListVPCs("key", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("expected error from provider")
	}
}

func TestListVPCs_EmptyCredentials(t *testing.T) {
	svc := NewVPCService()

	_, err := svc.ListVPCs("", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("expected error for empty accessKeyId")
	}

	_, err = svc.ListVPCs("key", "", "cn-hangzhou")
	if err == nil {
		t.Error("expected error for empty accessKeySecret")
	}

	_, err = svc.ListVPCs("key", "secret", "")
	if err == nil {
		t.Error("expected error for empty region")
	}
}

func TestNewVPCService(t *testing.T) {
	svc := NewVPCService()
	if svc == nil {
		t.Error("expected non-nil VPCService")
	}
}

func TestGetVPCDetail_MockProvider(t *testing.T) {
	mock := &provider.MockVPCProvider{
		VPCDetail: &provider.VPCDetail{
			VpcId:      "vpc-001",
			VpcName:    "test-vpc",
			CidrBlock:  "10.0.0.0/16",
			Status:     "Available",
			RegionId:   "cn-hangzhou",
			VSwitchIds: []string{"vsw-001", "vsw-002"},
		},
	}

	svc := NewVPCServiceWithProvider(func(a, b, c string) (provider.VPCProvider, error) {
		return mock, nil
	})

	result, err := svc.GetVPCDetail("key", "secret", "cn-hangzhou", "vpc-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result.VpcId != "vpc-001" {
		t.Errorf("expected VPC ID 'vpc-001', got '%s'", result.VpcId)
	}
	if len(result.VSwitchIds) != 2 {
		t.Errorf("expected 2 VSwitch IDs, got %d", len(result.VSwitchIds))
	}
}

func TestGetVPCDetail_EmptyVpcId(t *testing.T) {
	svc := NewVPCService()

	_, err := svc.GetVPCDetail("key", "secret", "cn-hangzhou", "")
	if err == nil {
		t.Error("expected error for empty vpcId")
	}
}

func TestListVSwitches_MockProvider(t *testing.T) {
	mock := &provider.MockVPCProvider{
		VSwitches: []provider.VSwitch{
			{VSwitchId: "vsw-001", VSwitchName: "test-vsw-1", CidrBlock: "10.0.1.0/24", ZoneId: "cn-hangzhou-a"},
			{VSwitchId: "vsw-002", VSwitchName: "test-vsw-2", CidrBlock: "10.0.2.0/24", ZoneId: "cn-hangzhou-b"},
		},
	}

	svc := NewVPCServiceWithProvider(func(a, b, c string) (provider.VPCProvider, error) {
		return mock, nil
	})

	result, err := svc.ListVSwitches("key", "secret", "cn-hangzhou", "vpc-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.VSwitches) != 2 {
		t.Errorf("expected 2 VSwitches, got %d", len(result.VSwitches))
	}
	if result.VSwitches[0].VSwitchId != "vsw-001" {
		t.Errorf("expected first VSwitch ID 'vsw-001', got '%s'", result.VSwitches[0].VSwitchId)
	}
}
