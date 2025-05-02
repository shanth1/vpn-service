package bubbletea

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shanth1/vpn-service/internal/core/ports"
)

type handler struct {
	service ports.PrimaryPort
}

func NewTUIHandler(service ports.PrimaryPort) *handler {
	return &handler{service: service}
}

func (h *handler) MustRun(ctx context.Context) {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{
		tabs:      []string{"1", "2", "3"},
		activeTab: 0,
		width:     10,
		height:    10,
	}
}
