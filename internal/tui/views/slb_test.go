package views

import (
	"cloud-manage/service"
	"testing"
)

func TestNewSLBView(t *testing.T) {
	view := NewSLBView()

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

func TestSLBView_SetSize(t *testing.T) {
	view := NewSLBView()
	view.SetSize(100, 50)

	if view.width != 100 {
		t.Errorf("expected width 100, got %d", view.width)
	}
	if view.height != 50 {
		t.Errorf("expected height 50, got %d", view.height)
	}
}

func TestSLBView_Update_Loaded(t *testing.T) {
	view := NewSLBView()
	view.SetSize(100, 50)

	msg := SLBLoaded{
		Result: &service.ListSLBsResult{
			SLBs: []service.SLBAdapter{
				{LoadBalancerId: "lb-001", LoadBalancerName: "test-slb-1", Address: "1.2.3.4"},
				{LoadBalancerId: "lb-002", LoadBalancerName: "test-slb-2", Address: "5.6.7.8"},
			},
		},
	}

	view.Update(msg)

	if view.loading {
		t.Error("expected loading to be false after update")
	}
}

func TestSLBView_Update_Error(t *testing.T) {
	view := NewSLBView()

	msg := SLBLoadError{
		Err: nil,
	}

	view.Update(msg)

	if view.loading {
		t.Error("expected loading to be false after error")
	}
}

func TestSLBView_Render_Loading(t *testing.T) {
	view := NewSLBView()
	view.loading = true

	result := view.Render()
	if result != "Loading..." {
		t.Errorf("expected 'Loading...', got '%s'", result)
	}
}

func TestSLBView_Render_Empty(t *testing.T) {
	view := NewSLBView()

	result := view.Render()
	if result == "" {
		t.Error("expected non-empty render")
	}
}
