package main

import (
	"context"
	"log"

	"github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea"
	"github.com/shanth1/vpn-service/internal/adapters/secondary/qr"
	"github.com/shanth1/vpn-service/internal/adapters/secondary/system"
	"github.com/shanth1/vpn-service/internal/adapters/secondary/wireguard"
	"github.com/shanth1/vpn-service/internal/common"
	"github.com/shanth1/vpn-service/internal/config"
	"github.com/shanth1/vpn-service/internal/core/ports"
	"github.com/shanth1/vpn-service/internal/core/usecase"
	"github.com/shanth1/vpn-service/pkg/configutil.go"
)

func main() {
	ctx := context.Background()

	cfg := config.TUI{}
	configutil.Load(configutil.GetConfigPath(), &cfg)

	var systemInfra ports.System
	if cfg.Env == common.EnvProd {
		systemInfra = system.NewInfra()
	} else {
		systemInfra = system.NewFakeInfra()
	}
	ifaceName, publicIP, err := systemInfra.GetNetworkInfo()
	if err != nil {
		log.Fatalf("network info: %v", err)
	}
	wgInfra := wireguard.NewInfra(publicIP, ifaceName)
	qrInfra := qr.NewInfra()

	service := usecase.NewVPNService(systemInfra, wgInfra, qrInfra)
	tuiHandler := bubbletea.NewTUIHandler(service)
	tuiHandler.MustRun(ctx)
}
