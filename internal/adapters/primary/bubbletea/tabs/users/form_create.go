package users

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tuicommon "github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/common"
	"github.com/shanth1/vpn-service/internal/core/domain"
)

func (m *Model) handleCreateFormKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch keypress := msg.String(); keypress {
	case "esc":
		m.state = stateListView
		m.err = nil
		m.successMsg = ""
		for i := range m.formInputs {
			m.formInputs[i].Blur()
		}
		return *m, nil

	case "enter":
		if m.formFocusIndex == len(m.formInputs) {
			tgUsername := strings.TrimSpace(m.formInputs[0].Value())
			name := strings.TrimSpace(m.formInputs[1].Value())
			if tgUsername == "" || name == "" {
				m.err = errors.New("Telegram Username and Name are required")
				return *m, nil
			}

			newUser := &domain.User{
				TG:      tgUsername,
				Name:    name,
				Surname: strings.TrimSpace(m.formInputs[2].Value()),
				Email:   strings.TrimSpace(m.formInputs[3].Value()),
			}

			m.state = stateSubmittingForm
			cmds = append(cmds, m.spinner.Tick, createUserCmd(m.core, newUser))
			return *m, tea.Batch(cmds...)
		}
		fallthrough

	case "tab", "shift+tab", "up", "down":
		s := msg.String()

		if s == "up" || s == "shift+tab" {
			m.formFocusIndex--
		} else {
			m.formFocusIndex++
		}

		if m.formFocusIndex > len(m.formInputs) {
			m.formFocusIndex = 0
		} else if m.formFocusIndex < 0 {
			m.formFocusIndex = len(m.formInputs)
		}

		for i := 0; i <= len(m.formInputs)-1; i++ {
			if i == m.formFocusIndex {
				cmds = append(cmds, m.formInputs[i].Focus())
				m.formInputs[i].PromptStyle = lipgloss.NewStyle().Foreground(tuicommon.AccentColor)
				m.formInputs[i].TextStyle = lipgloss.NewStyle()
				continue
			}
			m.formInputs[i].Blur()
			m.formInputs[i].PromptStyle = lipgloss.NewStyle()
			m.formInputs[i].TextStyle = lipgloss.NewStyle()
		}
		return *m, tea.Batch(cmds...)

	default:
		if m.formFocusIndex < len(m.formInputs) {
			var cmd tea.Cmd
			m.formInputs[m.formFocusIndex], cmd = m.formInputs[m.formFocusIndex].Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return *m, tea.Batch(cmds...)
}

func (m Model) renderCreateFormView() string {
	var b strings.Builder
	b.WriteString("Create New User ([Shift]Tab / Arrows to navigate, Enter on Submit, Esc to cancel)\n\n")

	for i := range m.formInputs {
		b.WriteString(m.formInputs[i].View() + "\n")
	}

	submitText := "Submit"
	if m.formFocusIndex == len(m.formInputs) {
		submitText = lipgloss.NewStyle().Foreground(tuicommon.AccentColor).Render("[ " + submitText + " ]")
	} else {
		submitText = "[ " + submitText + " ]"
	}
	b.WriteString("\n" + submitText + "\n")

	if m.err != nil {
		b.WriteString("\n" + tuicommon.ErrorStyle.Render(m.err.Error()))
	}

	formWidth := 0

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().Width(formWidth).Render(b.String()),
	)
}
