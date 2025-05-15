package main

import (
	"context"

	"github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea"
	"github.com/shanth1/vpn-service/internal/adapters/secondary/qr"
	"github.com/shanth1/vpn-service/internal/core/usecase"
)

func main() {
	ctx := context.Background()

	// TODO: protocol
	qrInfra := qr.NewQRInfra()

	qrInfra.Generate(ctx, "test", []byte{1, 2, 3, 4})
	service := usecase.NewVPNService(nil, qrInfra)
	tuiHandler := bubbletea.NewTUIHandler(service)
	tuiHandler.MustRun(ctx)

}
