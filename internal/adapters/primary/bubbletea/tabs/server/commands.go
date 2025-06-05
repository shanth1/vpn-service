package server

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	taskInstallID   = 0
	taskUninstallID = 1
	taskInitID      = 2
	taskCleanUpID   = 3
)

type taskResultMsg struct {
	id  int
	err error
}

func (m Model) processTask(taskID int) tea.Cmd {
	return func() tea.Msg {
		var err error
		switch taskID {
		case taskInstallID:
			err = m.core.Init(context.Background())
		case taskUninstallID:
			err = m.core.Uninstall(context.Background())
		case taskInitID:
			err = m.core.Init(context.Background())
		case taskCleanUpID:
			err = m.core.CleanUp(context.Background())
		default:
			err = errors.New("Unknown task")
		}

		return taskResultMsg{id: taskID, err: err}
	}
}
