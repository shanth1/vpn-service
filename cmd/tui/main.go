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
	"github.com/shanth1/vpn-service/pkg/configutil"
)

func main() {
	ctx := context.Background()

	cfg := config.TUI{}
	if err := configutil.Load(configutil.GetConfigPath(), &cfg); err != nil {
		log.Fatalf("load config: %w", err)
	}

	systemAdapter := newSystemAdapter(cfg.Env)
	ifaceName, publicIP, err := systemAdapter.GetNetworkInfo()
	if err != nil {
		log.Fatalf("network info: %v", err)
	}
	wgAdapter := newProtocolAdapter(cfg.Env, publicIP, ifaceName)
	qrAdapter := qr.NewAdapter()

	core := usecase.NewVPNCore(systemAdapter, wgAdapter, qrAdapter)
	tuiHandler := bubbletea.NewTUIHandler(core)
	tuiHandler.MustRun(ctx)
}

func newSystemAdapter(env string) ports.System {
	if env == common.EnvProd {
		return system.NewAdapter()
	}
	return system.NewFakeAdapter()
}

func newProtocolAdapter(env, serverPublicIP, netPublicIfaceName string) ports.Protocol {
	if env == common.EnvProd {
		return wireguard.NewAdapter(serverPublicIP, netPublicIfaceName)
	}
	return wireguard.NewFakeAdapter()
}
