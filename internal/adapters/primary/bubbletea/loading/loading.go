package loading

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tuicommon "github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/common"
)

const TabTitle = "Loading"

type contentLoadedMsg struct {
	content string
}

func loadContentCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		return contentLoadedMsg{content: "Контент успешно загружен.\n\nЭто демонстрация асинхронной загрузки данных."}
	}
}

type Model struct {
	width     int
	height    int
	isLoading bool
	spinner   spinner.Model
	content   string
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(tuicommon.AccentColor)
	return Model{
		isLoading: true,
		spinner:   s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(loadContentCmd(), m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case contentLoadedMsg:
		m.isLoading = false
		m.content = msg.content
		return m, nil

	case spinner.TickMsg:
		if m.isLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	}

	var spinCmd tea.Cmd
	if m.isLoading {
		m.spinner, spinCmd = m.spinner.Update(msg)
	}

	return m, spinCmd
}

func (m Model) View() string {
	if m.isLoading {
		loadingText := fmt.Sprintf("%s Загрузка данных, пожалуйста подождите...", m.spinner.View())
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			loadingText,
		)
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.content,
	)
}
