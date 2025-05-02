package bubbletea

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	activeTab int
	tabs      []string
	width     int
	height    int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			num, _ := strconv.Atoi(key)
			if num > len(m.tabs) {
				num = len(m.tabs)
			}
			m.activeTab = num - 1
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
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
			renderedTabs = append(renderedTabs, activeTabStyle.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(tab))
		}
	}

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	content := "Выбран таб: " + m.tabs[m.activeTab]

	return borderStyle.Render(tabsRow + lipgloss.NewStyle().Padding(2, 2).Render(content))
}
