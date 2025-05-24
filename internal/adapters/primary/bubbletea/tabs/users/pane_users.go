package users

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tuicommon "github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/common"
	"github.com/shanth1/vpn-service/internal/core/domain"
)

type userListItem struct {
	user *domain.User
}

func (uli userListItem) display() string {
	return fmt.Sprintf("%s (%s %s)", uli.user.TG, uli.user.Name, uli.user.Surname)
}

func (m Model) renderUserListPane(paneWidth, paneHeight int) string {
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(tuicommon.AccentColor).
		Padding(1, 2).Width(paneWidth).Height(paneHeight)

	if len(m.users) == 0 && m.err == nil { // No users and no error yet (might still be loading initial)
		return listStyle.Render("No users found. Press 'c' to create one.")
	}
	if len(m.users) == 0 && m.err != nil { // No users and an error occurred
		return listStyle.Render(fmt.Sprintf("Error loading users:\n%s\nPress 'c' to create or try reloading.", m.err.Error()))
	}

	var items []string
	items = append(items, "Users (j/k, Enter, 'c' to create):")
	items = append(items, "") // Spacer

	maxDisplayItems := paneHeight - 4 // Account for title, spacer, and borders/padding
	if maxDisplayItems < 1 {
		maxDisplayItems = 1
	}

	start := 0
	end := len(m.users)

	// Basic scrolling logic if list is too long for the pane
	if len(m.users) > maxDisplayItems {
		if m.selectedUserIndex >= maxDisplayItems/2 {
			start = m.selectedUserIndex - maxDisplayItems/2
			if start+maxDisplayItems > len(m.users) {
				start = len(m.users) - maxDisplayItems
			}
		}
		if start < 0 {
			start = 0
		}
		end = start + maxDisplayItems
		if end > len(m.users) {
			end = len(m.users)
		}
	}

	for i := start; i < end; i++ {
		user := m.users[i]
		itemStr := user.display()
		if i == m.selectedUserIndex {
			itemStr = lipgloss.NewStyle().Foreground(tuicommon.AccentColor).Bold(true).Render("> " + itemStr)
		} else {
			itemStr = "  " + itemStr
		}
		items = append(items, itemStr)
	}
	if len(m.users) > maxDisplayItems && end < len(m.users) {
		items = append(items, "  ...")
	}
	if len(m.users) > maxDisplayItems && start > 0 {
		items = append(items, "  ...") // Indicate items above
		// This part is tricky without a proper list component, might need reverse order or different scroll markers
	}

	return listStyle.Render(strings.Join(items, "\n"))
}
