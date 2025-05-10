package main

import (
	"context"

	"github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea"
	"github.com/shanth1/vpn-service/internal/core/usecase"
)

func main() {
	ctx := context.Background()

	service := usecase.NewVPNService()
	tuiHandler := bubbletea.NewTUIHandler(service)
	tuiHandler.MustRun(ctx)
}
