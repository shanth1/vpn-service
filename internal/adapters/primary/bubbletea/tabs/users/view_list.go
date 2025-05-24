package users

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) handleListViewKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "j", "down":
		if len(m.users) > 0 {
			m.selectedUserIndex = (m.selectedUserIndex + 1) % len(m.users)
		}
	case "k", "up":
		if len(m.users) > 0 {
			m.selectedUserIndex = (m.selectedUserIndex - 1 + len(m.users)) % len(m.users)
		}
	case "enter":
		if m.selectedUserIndex >= 0 && m.selectedUserIndex < len(m.users) {
			m.currentUserDetails = m.users[m.selectedUserIndex].user
			m.state = stateDetailsLoading
			return tea.Batch(m.spinner.Tick, fetchUserDetailsCmd(m.core, m.currentUserDetails.PublicKey))
		}
	case "c":
		m.state = stateCreateForm
		m.formFocusIndex = 0
		for i := range m.formInputs {
			m.formInputs[i].SetValue("")
			if i == m.formFocusIndex {
				m.formInputs[i].Focus()
			} else {
				m.formInputs[i].Blur()
			}
		}
		return textinput.Blink
	}
	return nil
}

func (m Model) renderListView() string {
	paneWidth := m.width / 2
	if paneWidth == 0 {
		paneWidth = 10
	} // Minimum width
	if m.width < 20 {
		paneWidth = m.width
	} // Full width if too small

	leftPane := m.renderUserListPane(paneWidth-1, m.height-1) // -1 for border

	rightPaneStyle := lipgloss.NewStyle().
		Padding(1, 2).Width(m.width - paneWidth - 1).Height(m.height - 1)

	var rightPaneContent string
	if m.selectedUserIndex != -1 && len(m.users) > 0 {
		rightPaneContent = "Select a user to see details (Enter)\n\nPress 'c' to create a new user."
	} else if len(m.users) == 0 && m.err == nil {
		rightPaneContent = "No users loaded.\n\nPress 'c' to create a new user."
	} else if m.err != nil {
		rightPaneContent = fmt.Sprintf("Error: %s\n\nPress 'c' to create or try reloading tab.", m.err.Error())
	} else {
		rightPaneContent = "Press 'c' to create a new user."
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftPane,
		rightPaneStyle.Render(rightPaneContent),
	)
}
