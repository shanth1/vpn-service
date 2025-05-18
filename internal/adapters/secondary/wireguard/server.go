package wireguard

import (
	"context"
	"fmt"
	"os"

	"github.com/shanth1/vpn-service/internal/common"
)

func (w *wireguardInfra) SetUpServer(ctx context.Context) error {
	if err := os.MkdirAll(w.serverDir, 0700); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	privateKeyBytes, publicKeyBytes, err := genWgKeys(ctx)
	if err != nil {
		return fmt.Errorf("gen keys: %w", err)
	}

	if err := os.WriteFile(w.serverPrivKeyPath, privateKeyBytes, 0600); err != nil {
		return fmt.Errorf("write server private key: %w", err)
	}

	if err := os.WriteFile(w.serverPrivKeyPath, publicKeyBytes, 0600); err != nil {
		return fmt.Errorf("write server public key: %w", w.serverPrivKeyPath, err)
	}

	serverConfContent := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s
ListenPort = %d
SaveConfig = true

PostUp = iptables -A FORWARD -i %%i -j ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE; iptables -A FORWARD -o %%i -j ACCEPT
PostDown = iptables -D FORWARD -i %%i -j ACCEPT; iptables -t nat -D POSTROUTING -o %s -j MASQUERADE; iptables -D FORWARD -o %%i -j ACCEPT
`, string(privateKeyBytes), w.serverVPNAddrCIDR, w.serverListenPort, w.publicNetIface, w.publicNetIface)

	if err := os.WriteFile(w.serverConfigPath, []byte(serverConfContent), 0600); err != nil {
		return fmt.Errorf("server configuration: %w", err)
	}

	return nil
}

func (w *wireguardInfra) TearDownServer(ctx context.Context) error {
	if err := w.StopService(ctx); err != nil {
		return fmt.Errorf("stop service: %w", w.serverDir, err)
	}

	_, _ = common.RunCommand(ctx, "systemctl", "disable", w.serviceName())

	if err := os.RemoveAll(w.serverDir); err != nil {
		return fmt.Errorf("remove directory: %w", w.serverDir, err)
	}

	return nil
}
