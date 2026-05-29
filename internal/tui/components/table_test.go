package components

import (
	"testing"
)

func TestNewTable(t *testing.T) {
	headers := []string{"ID", "Name", "Status"}
	table := NewTable(headers)

	if len(table.Headers) != 3 {
		t.Errorf("expected 3 headers, got %d", len(table.Headers))
	}
	if table.Headers[0] != "ID" {
		t.Errorf("expected first header 'ID', got '%s'", table.Headers[0])
	}
}

func TestTable_SetRows(t *testing.T) {
	table := NewTable([]string{"ID", "Name"})
	rows := [][]string{
		{"1", "test-1"},
		{"2", "test-2"},
		{"3", "test-3"},
	}

	table.SetRows(rows)

	if len(table.Rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(table.Rows))
	}
	if table.Cursor != 0 {
		t.Errorf("expected cursor 0, got %d", table.Cursor)
	}
}

func TestTable_SetRows_AdjustsCursor(t *testing.T) {
	table := NewTable([]string{"ID"})
	table.SetRows([][]string{{"1"}, {"2"}, {"3"}})
	table.Cursor = 5

	table.SetRows([][]string{{"1"}})

	if table.Cursor != 0 {
		t.Errorf("expected cursor 0, got %d", table.Cursor)
	}
}

func TestTable_MoveUp(t *testing.T) {
	table := NewTable([]string{"ID"})
	table.SetRows([][]string{{"1"}, {"2"}, {"3"}})
	table.Cursor = 2

	table.MoveUp()

	if table.Cursor != 1 {
		t.Errorf("expected cursor 1, got %d", table.Cursor)
	}
}

func TestTable_MoveUp_AtTop(t *testing.T) {
	table := NewTable([]string{"ID"})
	table.SetRows([][]string{{"1"}, {"2"}})
	table.Cursor = 0

	table.MoveUp()

	if table.Cursor != 0 {
		t.Errorf("expected cursor 0, got %d", table.Cursor)
	}
}

func TestTable_MoveDown(t *testing.T) {
	table := NewTable([]string{"ID"})
	table.SetRows([][]string{{"1"}, {"2"}, {"3"}})
	table.SetSize(80, 20)
	table.Cursor = 0

	table.MoveDown()

	if table.Cursor != 1 {
		t.Errorf("expected cursor 1, got %d", table.Cursor)
	}
}

func TestTable_MoveDown_AtBottom(t *testing.T) {
	table := NewTable([]string{"ID"})
	table.SetRows([][]string{{"1"}, {"2"}})
	table.Cursor = 1

	table.MoveDown()

	if table.Cursor != 1 {
		t.Errorf("expected cursor 1, got %d", table.Cursor)
	}
}

func TestTable_SelectedRow(t *testing.T) {
	table := NewTable([]string{"ID", "Name"})
	table.SetRows([][]string{
		{"1", "test-1"},
		{"2", "test-2"},
	})

	row := table.SelectedRow()
	if row == nil {
		t.Fatal("expected non-nil row")
	}
	if row[0] != "1" {
		t.Errorf("expected '1', got '%s'", row[0])
	}
}

func TestTable_SelectedRow_Empty(t *testing.T) {
	table := NewTable([]string{"ID"})

	row := table.SelectedRow()
	if row != nil {
		t.Error("expected nil row for empty table")
	}
}

func TestTable_Render_Empty(t *testing.T) {
	table := NewTable([]string{})

	result := table.Render()
	if result != "" {
		t.Error("expected empty render for no headers")
	}
}

func TestTable_Render_WithData(t *testing.T) {
	table := NewTable([]string{"ID", "Name"})
	table.SetRows([][]string{
		{"1", "test-1"},
		{"2", "test-2"},
	})
	table.SetSize(80, 20)

	result := table.Render()
	if result == "" {
		t.Error("expected non-empty render")
	}
}
