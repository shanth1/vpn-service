package bubbletea

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab struct {
	title   string
	content tea.Model
}

func newTab(title string, content tea.Model) *tab {
	return &tab{
		title:   title,
		content: content,
	}
}

type model struct {
	activeTab int
	tabs      []*tab
	width     int
	height    int
}

func (m model) Init() tea.Cmd {
	if len(m.tabs) > 0 && m.tabs[m.activeTab].content != nil {
		return m.tabs[m.activeTab].content.Init()
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			num, _ := strconv.Atoi(key)
			if num > 0 && num <= len(m.tabs) {
				m.activeTab = num - 1
				if m.tabs[m.activeTab].content != nil {
					cmd = m.tabs[m.activeTab].content.Init()
					cmds = append(cmds, cmd)
				}
			}
		default:
			if len(m.tabs) > 0 && m.tabs[m.activeTab].content != nil {
				var updatedContent tea.Model
				updatedContent, cmd = m.tabs[m.activeTab].content.Update(msg)
				m.tabs[m.activeTab].content = updatedContent
				cmds = append(cmds, cmd)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		tabsRowHeight := lipgloss.Height(m.renderTabsRow())
		contentWidth := m.width - 2
		contentHeight := m.height - 2 - tabsRowHeight
		contentSizeMsg := tea.WindowSizeMsg{
			Width:  contentWidth,
			Height: contentHeight,
		}
		for i := range m.tabs {
			if m.tabs[i].content != nil {
				var updatedContent tea.Model
				updatedContent, cmd = m.tabs[i].content.Update(contentSizeMsg)
				m.tabs[i].content = updatedContent
				cmds = append(cmds, cmd)
			}
		}
	default:
		if len(m.tabs) > 0 && m.tabs[m.activeTab].content != nil {
			var updatedContent tea.Model
			updatedContent, cmd = m.tabs[m.activeTab].content.Update(msg)
			m.tabs[m.activeTab].content = updatedContent
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m model) renderTabsRow() string {
	if m.width == 0 {
		return ""
	}
	tabStyle := lipgloss.NewStyle().Padding(0, 2)
	activeTabStyle := tabStyle.
		Background(lipgloss.Color("#fa0")).
		Foreground(lipgloss.Color("#000")).
		Bold(true)

	var renderedTabs []string
	for i, t := range m.tabs {
		title := getNumTabTitle(i, t.title)
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(title))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(title))
		}
	}
	return lipgloss.NewStyle().Width(m.width - 2).Render(lipgloss.JoinHorizontal(lipgloss.Left, renderedTabs...))
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	tabsRow := m.renderTabsRow()
	tabsRowHeight := lipgloss.Height(tabsRow)

	contentHeight := m.height - 2 - tabsRowHeight
	if contentHeight < 0 {
		contentHeight = 0
	}

	var contentStr string
	if len(m.tabs) > 0 && m.activeTab < len(m.tabs) && m.tabs[m.activeTab].content != nil {
		contentStr = m.tabs[m.activeTab].content.View()
	} else {
		contentStr = "No active content"
	}

	finalContent := lipgloss.JoinVertical(lipgloss.Left,
		tabsRow,
		lipgloss.NewStyle().Height(contentHeight).Render(contentStr),
	)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(m.width - 2).
		Height(m.height - 2).
		BorderForeground(lipgloss.Color("#fa0"))

	return borderStyle.Render(finalContent)
}

func getNumTabTitle(i int, title string) string {
	return fmt.Sprintf("(%d) %s", i+1, title)
}
