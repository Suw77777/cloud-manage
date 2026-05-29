package views

import (
	"cloud-manage/service"
	"testing"
)

func TestNewVPCView(t *testing.T) {
	view := NewVPCView()

	if view.loading {
		t.Error("expected loading to be false initially")
	}
	if view.showDetail {
		t.Error("expected showDetail to be false initially")
	}
	if view.service == nil {
		t.Error("expected non-nil service")
	}
}

func TestVPCView_SetSize(t *testing.T) {
	view := NewVPCView()
	view.SetSize(100, 50)

	if view.width != 100 {
		t.Errorf("expected width 100, got %d", view.width)
	}
	if view.height != 50 {
		t.Errorf("expected height 50, got %d", view.height)
	}
}

func TestVPCView_Update_Loaded(t *testing.T) {
	view := NewVPCView()
	view.SetSize(100, 50)

	msg := VPCLoaded{
		Result: &service.ListVPCsResult{
			VPCs: []service.VPCAdapter{
				{VpcId: "vpc-001", VpcName: "test-vpc-1", CidrBlock: "10.0.0.0/16"},
				{VpcId: "vpc-002", VpcName: "test-vpc-2", CidrBlock: "172.16.0.0/12"},
			},
		},
	}

	view.Update(msg)

	if view.loading {
		t.Error("expected loading to be false after update")
	}
}

func TestVPCView_Update_Error(t *testing.T) {
	view := NewVPCView()

	msg := VPCLoadError{
		Err: nil,
	}

	view.Update(msg)

	if view.loading {
		t.Error("expected loading to be false after error")
	}
}

func TestVPCView_Render_Loading(t *testing.T) {
	view := NewVPCView()
	view.loading = true

	result := view.Render()
	if result != "Loading..." {
		t.Errorf("expected 'Loading...', got '%s'", result)
	}
}

func TestVPCView_Render_Empty(t *testing.T) {
	view := NewVPCView()

	result := view.Render()
	if result == "" {
		t.Error("expected non-empty render")
	}
}
