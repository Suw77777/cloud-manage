package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Table struct {
	Headers []string
	Rows    [][]string
	Cursor  int
	Offset  int
	width   int
	height  int
}

func NewTable(headers []string) Table {
	return Table{
		Headers: headers,
	}
}

func (t *Table) SetSize(width, height int) {
	t.width = width
	t.height = height
}

func (t *Table) SetRows(rows [][]string) {
	t.Rows = rows
	if t.Cursor >= len(rows) {
		t.Cursor = len(rows) - 1
	}
	if t.Cursor < 0 {
		t.Cursor = 0
	}
}

func (t *Table) MoveUp() {
	if t.Cursor > 0 {
		t.Cursor--
		if t.Cursor < t.Offset {
			t.Offset = t.Cursor
		}
	}
}

func (t *Table) MoveDown() {
	if t.Cursor < len(t.Rows)-1 {
		t.Cursor++
		visibleRows := t.height - 3
		if t.Cursor >= t.Offset+visibleRows {
			t.Offset = t.Cursor - visibleRows + 1
		}
	}
}

func (t *Table) SelectedRow() []string {
	if t.Cursor >= 0 && t.Cursor < len(t.Rows) {
		return t.Rows[t.Cursor]
	}
	return nil
}

func (t Table) Render() string {
	if len(t.Headers) == 0 {
		return ""
	}

	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00BFFF")).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#666666"))

	headerStr := ""
	for i, h := range t.Headers {
		if i > 0 {
			headerStr += "  "
		}
		headerStr += fmt.Sprintf("%-20s", h)
	}
	b.WriteString(headerStyle.Render(headerStr))
	b.WriteString("\n")

	// Rows
	visibleRows := t.height - 3
	if visibleRows < 1 {
		visibleRows = 1
	}

	end := t.Offset + visibleRows
	if end > len(t.Rows) {
		end = len(t.Rows)
	}

	for i := t.Offset; i < end; i++ {
		row := t.Rows[i]
		rowStr := ""
		for j, cell := range row {
			if j > 0 {
				rowStr += "  "
			}
			rowStr += fmt.Sprintf("%-20s", cell)
		}

		if i == t.Cursor {
			selectedStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("#1A1A2E")).
				Foreground(lipgloss.Color("#FFFFFF"))
			b.WriteString(selectedStyle.Render(rowStr))
		} else {
			b.WriteString(rowStr)
		}
		b.WriteString("\n")
	}

	return b.String()
}
