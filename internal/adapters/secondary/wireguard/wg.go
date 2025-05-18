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

type wireguardInfra struct {
	serverPublicIP    string
	serverListenPort  int
	serverVPNAddrCIDR string
	serverDir         string
	clientsDir        string
	serverConfigPath  string
	serverPubKeyPath  string
	serverPrivKeyPath string
	publicNetIface    string
	wgInterfaceName   string
}

func NewWireguardInfra(serverPublicIP string) ports.Protocol {
	rootDir := "/etc/wireguard"
	wgInterface := "wg0"

	return &wireguardInfra{
		serverPublicIP:    serverPublicIP,
		serverListenPort:  51820,
		serverDir:         rootDir,
		clientsDir:        filepath.Join(rootDir, "clients"),
		serverConfigPath:  filepath.Join(rootDir, getConfFileName(wgInterface)),
		serverPubKeyPath:  getPublicKeyPath(rootDir),
		serverPrivKeyPath: getPrivateKeyPath(rootDir),
		publicNetIface:    "eth0",
		wgInterfaceName:   wgInterface,
		serverVPNAddrCIDR: "10.0.0.1/24",
	}
}
