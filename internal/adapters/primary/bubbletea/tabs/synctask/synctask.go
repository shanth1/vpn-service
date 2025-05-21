package synctask

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tuicommon "github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/common"
	"github.com/shanth1/vpn-service/internal/core/ports"
)

const TabTitle = "Sync Task"

type taskResultMsg struct {
	taskIndex int
	result    string
	err       error
}

type taskStartMsg struct {
	taskIndex int
}

type taskInfo struct {
	id   int
	name string
}

type Model struct {
	core             ports.PrimaryPort
	width            int
	height           int
	tasks            []taskInfo
	selectedTask     int
	runningTaskIndex int
	lastResult       string
	isLoading        bool
	err              error
}

func New(core ports.PrimaryPort) Model {
	tasks := []taskInfo{
		{id: 0, name: "Задача A (2 сек, успех/ошибка 50%)"},
		{id: 1, name: "Задача B (3 сек, всегда успех)"},
		{id: 2, name: "Задача C (1 сек, всегда ошибка)"},
	}
	return Model{
		core:             core,
		tasks:            tasks,
		selectedTask:     0,
		runningTaskIndex: -1,
		isLoading:        false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func simulateTask(taskID int) (string, error) {
	switch taskID {
	case 0:
		time.Sleep(2 * time.Second)
		if rand.Intn(2) == 0 {
			return fmt.Sprintf("Задача A (%d) успешно завершена!", taskID), nil
		}
		return "", fmt.Errorf("задача A (%d) завершилась с ошибкой", taskID)
	case 1:
		time.Sleep(3 * time.Second)
		return fmt.Sprintf("Задача B (%d) всегда успешна.", taskID), nil
	case 2:
		time.Sleep(1 * time.Second)
		return "", fmt.Errorf("задача C (%d) всегда завершается ошибкой", taskID)
	default:
		return "", errors.New("неизвестная задача")
	}
}

func runTaskCmd(taskIndex int, taskID int) tea.Cmd {
	return func() tea.Msg {
		result, err := simulateTask(taskID)

		return taskResultMsg{taskIndex: taskIndex, result: result, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.isLoading {
			return m, nil
		}
		switch msg.String() {
		case "j", "down":
			m.selectedTask = (m.selectedTask + 1) % len(m.tasks)
		case "k", "up":
			m.selectedTask = (m.selectedTask - 1 + len(m.tasks)) % len(m.tasks)
		case "enter":
			if m.runningTaskIndex == -1 {
				m.isLoading = true
				m.runningTaskIndex = m.selectedTask
				m.lastResult = ""
				m.err = nil
				selectedTaskInfo := m.tasks[m.selectedTask]
				return m, runTaskCmd(m.selectedTask, selectedTaskInfo.id)
			}
		}

	case taskResultMsg:
		m.isLoading = false
		m.runningTaskIndex = -1
		if msg.err != nil {
			m.err = msg.err
			m.lastResult = ""
		} else {
			m.err = nil
			m.lastResult = msg.result
		}
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	listStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(tuicommon.AccentColor).Padding(1, 2)
	statusStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(tuicommon.AccentColor).Padding(1, 2)
	selectedItemStyle := lipgloss.NewStyle().Foreground(tuicommon.AccentColor).Bold(true)
	loadingStyle := lipgloss.NewStyle().Foreground(tuicommon.AccentColor).Italic(true)

	var listBuilder strings.Builder
	listBuilder.WriteString("Выберите задачу (Up/Down) и нажмите Enter:\n\n")
	for i, task := range m.tasks {
		cursor := "  "
		style := lipgloss.NewStyle()
		if i == m.selectedTask {
			cursor = "> "
			style = selectedItemStyle
		}
		runningIndicator := ""
		if i == m.runningTaskIndex {
			runningIndicator = loadingStyle.Render(" (выполняется...)")
		}
		listBuilder.WriteString(fmt.Sprintf("%s%s%s\n", style.Render(cursor), style.Render(task.name), runningIndicator))
	}

	var statusBuilder strings.Builder
	statusBuilder.WriteString("Статус:\n\n")
	if m.isLoading && m.runningTaskIndex != -1 {
		statusBuilder.WriteString(loadingStyle.Render(fmt.Sprintf("Выполняется: %s", m.tasks[m.runningTaskIndex].name)))
	} else if m.err != nil {
		statusBuilder.WriteString(tuicommon.ErrorStyle.Render(fmt.Sprintf("Ошибка: %v", m.err)))
	} else if m.lastResult != "" {
		statusBuilder.WriteString(tuicommon.SuccessStyle.Render(m.lastResult))
	} else {
		statusBuilder.WriteString("Ожидание запуска задачи...")
	}

	availableWidth := m.width - listStyle.GetHorizontalFrameSize() - statusStyle.GetHorizontalFrameSize()
	listWidth := availableWidth / 2
	statusWidth := availableWidth - listWidth

	listStyle = listStyle.Width(listWidth).Height(m.height - listStyle.GetVerticalFrameSize() - 2)
	statusStyle = statusStyle.Width(statusWidth).Height(m.height - statusStyle.GetVerticalFrameSize() - 2)

	return lipgloss.JoinHorizontal(0,
		listStyle.Render(listBuilder.String()),
		statusStyle.Render(statusBuilder.String()),
	)
}
