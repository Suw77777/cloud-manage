package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Detail struct {
	Title  string
	Fields map[string]string
	width  int
	height int
}

func NewDetail(title string) Detail {
	return Detail{
		Title:  title,
		Fields: make(map[string]string),
	}
}

func (d *Detail) SetSize(width, height int) {
	d.width = width
	d.height = height
}

func (d *Detail) SetField(key, value string) {
	d.Fields[key] = value
}

func (d *Detail) Clear() {
	d.Fields = make(map[string]string)
}

func (d Detail) Render() string {
	if len(d.Fields) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("No details available")
	}

	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00BFFF"))

	b.WriteString(titleStyle.Render(d.Title))
	b.WriteString("\n\n")

	// Fields
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF")).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	for key, value := range d.Fields {
		b.WriteString(labelStyle.Render(key + ":"))
		b.WriteString(valueStyle.Render(value))
		b.WriteString("\n")
	}

	// Border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00BFFF")).
		Padding(1, 2).
		Width(d.width - 2)

	return borderStyle.Render(b.String())
}
