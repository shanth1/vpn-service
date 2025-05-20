package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/shanth1/vpn-service/internal/adapters/primary/bubbletea"
	"github.com/shanth1/vpn-service/internal/adapters/secondary/qr"
	"github.com/shanth1/vpn-service/internal/adapters/secondary/system"
	"github.com/shanth1/vpn-service/internal/adapters/secondary/wireguard"
	"github.com/shanth1/vpn-service/internal/config"
	"github.com/shanth1/vpn-service/internal/core/usecase"
	"github.com/shanth1/vpn-service/pkg/reader"
)

func main() {
	ctx := context.Background()

	configPath := flag.String("config", "", "Path to the config file")
	flag.Parse()
	if *configPath == "" {
		log.Fatal("usage: app -config /path/to/config.yaml")
	}

	cfg := config.TUI{}
	reader.Load(*configPath, &cfg)
	fmt.Println(cfg)

	systemInfra := system.NewInfra()
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
