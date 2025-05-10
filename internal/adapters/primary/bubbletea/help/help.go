package helppage

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tuicommon "github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/common"
)

const TabTitle = "Help"

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
	width    int
	height   int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerView := m.headerView()
		footerView := m.footerView()
		headerHeight := lipgloss.Height(headerView)
		footerHeight := lipgloss.Height(footerView)
		verticalMarginHeight := headerHeight + footerHeight

		viewportHeight := m.height - verticalMarginHeight
		if viewportHeight < 0 {
			viewportHeight = 0
		}

		if !m.ready {
			m.viewport = viewport.New(m.width, viewportHeight)
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = m.width
			m.viewport.Height = viewportHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready || m.width == 0 || m.height == 0 {
		return "Initializing help page..."
	}

	headerView := m.headerView()
	footerView := m.footerView()

	return lipgloss.JoinVertical(lipgloss.Left,
		headerView,
		m.viewport.View(),
		footerView,
	)
}

func (m model) headerView() string {
	title := tuicommon.TitleStyle.Render("Help Content Header")
	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(title)
}

func (m model) footerView() string {
	info := tuicommon.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Right).Render(info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func New() tea.Model {
	content := `# Help Section

This is the help page content. You can scroll up and down using the arrow keys, j/k, or page up/down.

## Subsection 1

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

> Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.

## Subsection 2

Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

 - List item 1
 - List item 2
   - Nested item
 - List item 3

More text to ensure scrolling is necessary.
Even more text.
Line after line.
...
...
...
...
...
...
...
Final line of help content.`

	return &model{content: content, ready: false}
}
