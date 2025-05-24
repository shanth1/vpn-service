package users

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tuicommon "github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/common"
)

func (m *Model) handleDetailsViewKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		m.currentUserDetails = nil
		m.currentUserTraffic = nil
		m.err = nil
		m.successMsg = ""
		m.state = stateListView
	case "D":
		if m.currentUserDetails != nil {
			m.state = stateSubmittingForm
			return tea.Batch(m.spinner.Tick, deleteUserCmd(m.core, m.currentUserDetails.TG))
		}
	}
	return nil
}

func (m Model) renderDetailsView() string {
	listPaneWidth := m.width / 3
	if listPaneWidth == 0 {
		listPaneWidth = 10
	}
	if m.width < 30 {
		listPaneWidth = m.width / 2
	}

	leftPane := m.renderUserListPane(listPaneWidth-1, m.height-1)

	detailsPaneStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, false).
		BorderForeground(tuicommon.AccentColor).
		Padding(1, 2).Width(m.width - listPaneWidth - 1).Height(m.height - 1)

	var detailsContent strings.Builder
	if m.currentUserDetails != nil {
		detailsContent.WriteString(fmt.Sprintf("Details for: %s %s (%s)\n",
			m.currentUserDetails.Name, m.currentUserDetails.Surname, m.currentUserDetails.TG))
		detailsContent.WriteString(fmt.Sprintf("ID: %s\n", m.currentUserDetails.TG))
		detailsContent.WriteString(fmt.Sprintf("Email: %s\n", m.currentUserDetails.Email))
		detailsContent.WriteString(fmt.Sprintf("Public Key: %s\n", m.currentUserDetails.PublicKey))
		detailsContent.WriteString("\n--- Traffic ---\n")
		if m.currentUserTraffic != nil {
			detailsContent.WriteString(fmt.Sprintf("Uploaded: %.2f MB\n", float64(m.currentUserTraffic.ReceivedBytes)/(1024*1024)))
			detailsContent.WriteString(fmt.Sprintf("Downloaded: %.2f MB\n", float64(m.currentUserTraffic.TransferBytes)/(1024*1024)))
		} else if m.err != nil && strings.Contains(m.err.Error(), "user details") {
			detailsContent.WriteString(tuicommon.ErrorStyle.Render(fmt.Sprintf("Could not load traffic: %v", m.err)))
		} else {
			detailsContent.WriteString("No traffic data available or still loading.")
		}

		detailsContent.WriteString("\n\n")
		detailsContent.WriteString("Press 'shidt + d' to Delete User\n")
		detailsContent.WriteString("Press 'esc' to go Back to list")

	} else if m.err != nil {
		detailsContent.WriteString(tuicommon.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		detailsContent.WriteString("\n\nPress 'esc' to go Back to list")
	} else {
		detailsContent.WriteString("No user selected or details not loaded.")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftPane,
		detailsPaneStyle.Render(detailsContent.String()),
	)
}
