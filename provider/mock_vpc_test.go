package provider

import (
	"errors"
	"testing"
)

func TestMockVPCProvider_ListVPCs_Success(t *testing.T) {
	mock := &MockVPCProvider{
		VPCs: []VPC{
			{VpcId: "vpc-001", VpcName: "test-vpc-1", CidrBlock: "10.0.0.0/16"},
			{VpcId: "vpc-002", VpcName: "test-vpc-2", CidrBlock: "172.16.0.0/12"},
		},
	}

	result, err := mock.ListVPCs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 VPCs, got %d", len(result))
	}
	if result[0].VpcId != "vpc-001" {
		t.Errorf("expected first VPC ID 'vpc-001', got '%s'", result[0].VpcId)
	}
}

func TestMockVPCProvider_ListVPCs_Error(t *testing.T) {
	mock := &MockVPCProvider{
		Err: errors.New("API error"),
	}

	_, err := mock.ListVPCs()
	if err == nil {
		t.Error("expected error")
	}
	if err.Error() != "API error" {
		t.Errorf("expected 'API error', got '%s'", err.Error())
	}
}

func TestMockVPCProvider_GetVPCDetail_Success(t *testing.T) {
	mock := &MockVPCProvider{
		VPCDetail: &VPCDetail{
			VpcId:      "vpc-001",
			VpcName:    "test-vpc",
			CidrBlock:  "10.0.0.0/16",
			VSwitchIds: []string{"vsw-001", "vsw-002"},
		},
	}

	result, err := mock.GetVPCDetail("vpc-001")
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

func TestMockVPCProvider_GetVPCDetail_Error(t *testing.T) {
	mock := &MockVPCProvider{
		Err: errors.New("not found"),
	}

	_, err := mock.GetVPCDetail("vpc-001")
	if err == nil {
		t.Error("expected error")
	}
}

func TestMockVPCProvider_ListVSwitches_Success(t *testing.T) {
	mock := &MockVPCProvider{
		VSwitches: []VSwitch{
			{VSwitchId: "vsw-001", VSwitchName: "test-vsw-1", CidrBlock: "10.0.1.0/24"},
			{VSwitchId: "vsw-002", VSwitchName: "test-vsw-2", CidrBlock: "10.0.2.0/24"},
		},
	}

	result, err := mock.ListVSwitches("vpc-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 VSwitches, got %d", len(result))
	}
}

func TestMockVPCProvider_ListVSwitches_Error(t *testing.T) {
	mock := &MockVPCProvider{
		Err: errors.New("API error"),
	}

	_, err := mock.ListVSwitches("vpc-001")
	if err == nil {
		t.Error("expected error")
	}
}
