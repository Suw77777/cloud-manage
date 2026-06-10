package views

import (
	"cloud-manage/internal/handler"
	"cloud-manage/service"
	"strings"
	"testing"
)

func TestNewECSView(t *testing.T) {
	view := NewECSView()

	if view.loading {
		t.Error("expected loading to be false initially")
	}
	if view.showDetail {
		t.Error("expected showDetail to be false initially")
	}
	if view.handler == nil {
		t.Error("expected non-nil handler")
	}
}

func TestECSView_SetSize(t *testing.T) {
	view := NewECSView()
	view.SetSize(100, 50)

	if view.width != 100 {
		t.Errorf("expected width 100, got %d", view.width)
	}
	if view.height != 50 {
		t.Errorf("expected height 50, got %d", view.height)
	}
}

func TestECSView_Update_Loaded(t *testing.T) {
	view := NewECSView()
	view.SetSize(100, 50)

	msg := ECSLoaded{
		Result: &handler.ECSListResult{
			Region: "cn-hangzhou",
			Instances: []service.ECSInstanceAdapter{
				{InstanceId: "i-001", InstanceName: "test-1", Status: "Running"},
				{InstanceId: "i-002", InstanceName: "test-2", Status: "Stopped"},
			},
			Total: 2,
		},
	}

	view.Update(msg)

	if view.loading {
		t.Error("expected loading to be false after update")
	}
}

func TestECSView_Update_Error(t *testing.T) {
	view := NewECSView()

	msg := ECSLoadError{
		Err: nil,
	}

	view.Update(msg)

	if view.loading {
		t.Error("expected loading to be false after error")
	}
}

func TestECSView_Render_Loading(t *testing.T) {
	view := NewECSView()
	view.loading = true
	view.region = "cn-hangzhou"

	result := view.Render()
	if !strings.Contains(result, "Loading") {
		t.Errorf("expected 'Loading' in render, got '%s'", result)
	}
}

func TestECSView_Render_Empty(t *testing.T) {
	view := NewECSView()

	result := view.Render()
	if result == "" {
		t.Error("expected non-empty render")
	}
}
