package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Search struct {
	input  textinput.Model
	active bool
	width  int
}

func NewSearch() Search {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 100
	ti.Width = 50
	return Search{
		input: ti,
	}
}

func (s *Search) SetWidth(width int) {
	s.width = width
	s.input.Width = width - 4
}

func (s *Search) Activate() {
	s.active = true
	s.input.Focus()
}

func (s *Search) Deactivate() {
	s.active = false
	s.input.Blur()
	s.input.SetValue("")
}

func (s *Search) IsActive() bool {
	return s.active
}

func (s *Search) Value() string {
	return s.input.Value()
}

func (s *Search) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return cmd
}

func (s Search) Render() string {
	if !s.active {
		return ""
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00BFFF")).
		Padding(0, 1).
		Width(s.width - 2)

	return style.Render(s.input.View())
}
