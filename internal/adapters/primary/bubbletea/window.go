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
		for i := range m.tabs {
			if m.tabs[i].content != nil {
				var updatedContent tea.Model
				updatedContent, cmd = m.tabs[i].content.Update(msg)
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

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(m.width - 2).
		Height(m.height - 2).
		BorderForeground(lipgloss.Color("#fa0"))

	tabStyle := lipgloss.NewStyle().
		Padding(0, 2)
	activeTabStyle := tabStyle.
		Background(lipgloss.Color("#fa0")).
		Foreground(lipgloss.Color("#000")).
		Bold(true)

	var renderedTabs []string
	for i, tab := range m.tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(getNumTabTitle(i, tab.title)))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(getNumTabTitle(i, tab.title)))
		}
	}

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	var content string
	if len(m.tabs) > 0 && m.tabs[m.activeTab].content != nil {
		content = m.tabs[m.activeTab].content.View()
	} else {
		content = "No active content"
	}

	return borderStyle.Render(tabsRow + lipgloss.NewStyle().Padding(2, 2).Render(content))
}

func getNumTabTitle(i int, title string) string {
	return fmt.Sprintf("(%d) %s", i+1, title)
}
