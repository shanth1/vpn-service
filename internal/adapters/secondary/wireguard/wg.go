package wireguard

import (
	"path/filepath"

	"github.com/shanth1/vpn-service/internal/core/ports"
)

/*
/etc/wireguard
├── wg0.conf
├── privatekey
├── publickey
└── clients
   └── cilentid
      ├── cilentid.conf
      ├── privatekey
      └── publickey
*/

type wgAdapter struct {
	serverPublicIP     string
	serverListenPort   int
	serverVPNAddrCIDR  string
	serverDir          string
	clientsDir         string
	serverConfigPath   string
	serverPubKeyPath   string
	serverPrivKeyPath  string
	netPublicIfaceName string
	wgInterfaceName    string
}

func NewAdapter(serverPublicIP, netPublicIfaceName string) ports.Protocol {
	rootDir := "/etc/wireguard"
	wgInterface := "wg0"

	return &wgAdapter{
		serverPublicIP:     serverPublicIP,
		serverListenPort:   51820,
		serverDir:          rootDir,
		clientsDir:         filepath.Join(rootDir, "clients"),
		serverConfigPath:   filepath.Join(rootDir, getConfFileName(wgInterface)),
		serverPubKeyPath:   getPublicKeyPath(rootDir),
		serverPrivKeyPath:  getPrivateKeyPath(rootDir),
		netPublicIfaceName: netPublicIfaceName,
		wgInterfaceName:    wgInterface,
		serverVPNAddrCIDR:  "10.0.0.1/24",
	}
}

type wgFakeAdapter struct{}

func NewFakeAdapter() ports.Protocol {
	return &wgFakeAdapter{}
}
