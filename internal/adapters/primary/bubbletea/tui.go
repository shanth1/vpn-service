package bubbletea

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	helppage "github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/tabs/help"
	"github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/tabs/loading"
	"github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea/tabs/synctask"
	"github.com/shanth1/vpn-service/internal/core/ports"
)

type handler struct {
	core ports.PrimaryPort
}

func NewTUIHandler(core ports.PrimaryPort) *handler {
	return &handler{core: core}
}

func (h *handler) MustRun(ctx context.Context) {
	p := tea.NewProgram(initialModel(h.core), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel(core ports.PrimaryPort) model {
	return model{
		tabs: []*tab{
			newTab(helppage.TabTitle, helppage.New()),
			newTab(synctask.TabTitle, synctask.New(core)),
			newTab(loading.TabTitle, loading.New()),
			{
				title:   "Service",
				content: nil,
			},
			{
				title:   "Monitoring",
				content: nil,
			},
			{
				title:   "Users",
				content: nil,
			},
		},
		activeTab: 0,
		width:     10,
		height:    10,
	}
}
