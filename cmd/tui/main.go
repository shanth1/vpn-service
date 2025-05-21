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

	var systemAdapter ports.System
	if cfg.Env == common.EnvProd {
		systemAdapter = system.NewAdapter()
	} else {
		systemAdapter = system.NewFakeAdapter()
	}
	ifaceName, publicIP, err := systemAdapter.GetNetworkInfo()
	if err != nil {
		log.Fatalf("network info: %v", err)
	}
	wgAdapter := wireguard.NewAdapter(publicIP, ifaceName)
	qrAdapter := qr.NewAdapter()

	core := usecase.NewVPNCore(systemAdapter, wgAdapter, qrAdapter)
	tuiHandler := bubbletea.NewTUIHandler(core)
	tuiHandler.MustRun(ctx)
}
