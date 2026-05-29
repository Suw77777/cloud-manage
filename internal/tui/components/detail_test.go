package components

import (
	"testing"
)

func TestNewDetail(t *testing.T) {
	detail := NewDetail("Test Detail")

	if detail.Title != "Test Detail" {
		t.Errorf("expected title 'Test Detail', got '%s'", detail.Title)
	}
	if detail.Fields == nil {
		t.Error("expected non-nil Fields map")
	}
}

func TestDetail_SetField(t *testing.T) {
	detail := NewDetail("Test")
	detail.SetField("Key1", "Value1")
	detail.SetField("Key2", "Value2")

	if len(detail.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(detail.Fields))
	}
	if detail.Fields["Key1"] != "Value1" {
		t.Errorf("expected 'Value1', got '%s'", detail.Fields["Key1"])
	}
}

func TestDetail_SetField_Overwrite(t *testing.T) {
	detail := NewDetail("Test")
	detail.SetField("Key", "Value1")
	detail.SetField("Key", "Value2")

	if detail.Fields["Key"] != "Value2" {
		t.Errorf("expected 'Value2', got '%s'", detail.Fields["Key"])
	}
}

func TestDetail_Clear(t *testing.T) {
	detail := NewDetail("Test")
	detail.SetField("Key1", "Value1")
	detail.SetField("Key2", "Value2")

	detail.Clear()

	if len(detail.Fields) != 0 {
		t.Errorf("expected 0 fields after clear, got %d", len(detail.Fields))
	}
}

func TestDetail_Render_Empty(t *testing.T) {
	detail := NewDetail("Test")
	detail.SetSize(80, 20)

	result := detail.Render()
	if result == "" {
		t.Error("expected non-empty render for empty detail")
	}
}

func TestDetail_Render_WithFields(t *testing.T) {
	detail := NewDetail("Test Detail")
	detail.SetField("ID", "test-001")
	detail.SetField("Name", "Test Name")
	detail.SetSize(80, 20)

	result := detail.Render()
	if result == "" {
		t.Error("expected non-empty render")
	}
}
