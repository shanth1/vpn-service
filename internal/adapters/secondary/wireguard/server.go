package wireguard

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/shanth1/vpn-service/internal/common"
)

func (w *wireguardInfra) SetUpServer(ctx context.Context) error {
	if err := os.MkdirAll(w.serverDir, 0700); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	privateKeyBytes, err := common.RunCommand(ctx, "wg", "genkey")
	if err != nil {
		return fmt.Errorf("wg genkey: %w", err)
	}
	privateKeyBytes = bytes.TrimSpace(privateKeyBytes)
	if err := os.WriteFile(w.serverPrivKeyPath, privateKeyBytes, 0600); err != nil {
		return fmt.Errorf("write server private key: %w", err)
	}

	cmdPubKey := exec.CommandContext(ctx, "wg", "pubkey")
	cmdPubKey.Stdin = bytes.NewReader(privateKeyBytes)
	publicKeyBytes, err := cmdPubKey.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("wg pubkey failed: %w. Stderr: %s", err, string(exitErr.Stderr))
		}
		return fmt.Errorf("wg pubkey failed: %w", err)
	}
	publicKeyBytes = bytes.TrimSpace(publicKeyBytes)
	if err := os.WriteFile(w.serverPrivKeyPath, publicKeyBytes, 0600); err != nil {
		return fmt.Errorf("write server public key: %w", w.serverPrivKeyPath, err)
	}
	privateKey := string(privateKeyBytes)

	serverConfContent := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s
ListenPort = %d
SaveConfig = true

PostUp = iptables -A FORWARD -i %%i -j ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE; iptables -A FORWARD -o %%i -j ACCEPT
PostDown = iptables -D FORWARD -i %%i -j ACCEPT; iptables -t nat -D POSTROUTING -o %s -j MASQUERADE; iptables -D FORWARD -o %%i -j ACCEPT
`, privateKey, w.serverVPNAddrCIDR, w.serverListenPort, w.publicNetIface, w.publicNetIface)

	if err := os.WriteFile(w.serverConfigPath, []byte(serverConfContent), 0600); err != nil {
		return fmt.Errorf("server configuration: %w", err)
	}

	return nil
}

func (w *wireguardInfra) TearDownServer(ctx context.Context) error {
	if err := w.StopService(ctx); err != nil {
		// TODO: log warn
	}

	if _, err := common.RunCommand(ctx, "systemctl", "disable", w.serviceName()); err != nil {
		// TODO: log warn
	}

	if err := os.RemoveAll(w.serverDir); err != nil {
		return fmt.Errorf("remove directory: %w", w.serverDir, err)
	}

	return nil
}
