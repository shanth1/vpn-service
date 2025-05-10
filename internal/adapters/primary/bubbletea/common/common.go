package tuicommon

import "github.com/charmbracelet/lipgloss"

const (
	AccentColor  = lipgloss.Color("#fa0")
	PrimaryColor = lipgloss.Color("000")

	// warningColor = lipgloss.Color("82")
	successColor = lipgloss.Color("82")
	errorColor   = lipgloss.Color("196")
)

var (
	SuccessStyle = lipgloss.NewStyle().Foreground(successColor)
	ErrorStyle   = lipgloss.NewStyle().Foreground(errorColor)
	TitleStyle   = lipgloss.NewStyle().Bold(true).Margin(0, 1)
	InfoStyle    = lipgloss.NewStyle().Margin(0, 1)
)
