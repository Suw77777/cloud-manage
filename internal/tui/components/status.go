package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Status struct {
	Message string
	Error   string
	Region  string
	width   int
}

func NewStatus() Status {
	return Status{
		Region: "cn-hangzhou",
	}
}

func (s *Status) SetWidth(width int) {
	s.width = width
}

func (s *Status) SetMessage(msg string) {
	s.Message = msg
	s.Error = ""
}

func (s *Status) SetError(err string) {
	s.Error = err
	s.Message = ""
}

func (s *Status) SetRegion(region string) {
	s.Region = region
}

func (s Status) Render() string {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Width(s.width)

	left := ""
	if s.Error != "" {
		left = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Render("Error: " + s.Error)
	} else if s.Message != "" {
		left = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF88")).
			Render(s.Message)
	} else {
		left = "Ready"
	}

	right := fmt.Sprintf("Region: %s | q: quit | /: search | r: refresh", s.Region)

	gap := s.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 0 {
		gap = 0
	}

	return style.Render(left + lipgloss.NewStyle().Width(gap).Render("") + right)
}
