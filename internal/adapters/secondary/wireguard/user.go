package wireguard

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/shanth1/vpn-service/internal/common"
	"github.com/shanth1/vpn-service/internal/core/domain"
)

func (w *wireguardInfra) AddUser(ctx context.Context, user *domain.User) error {
	if user.Name == "" {
		return errors.New("user name is required")
	}
	if w.serverPublicIP == "" {
		return errors.New("serverPublicIP is not configured")
	}

	clientDir := filepath.Join(w.clientsDir, user.TG)
	clientPrivKeyPath := filepath.Join(clientDir, "privatekey")
	clientPubKeyPath := filepath.Join(clientDir, "publickey")
	clientConfFileName := user.TG + ".conf"
	clientConfPath := filepath.Join(clientDir, clientConfFileName)

	log.Printf("Creating directory for client %s: %s", user.TG, clientDir)
	if err := os.MkdirAll(clientDir, 0700); err != nil {
		return fmt.Errorf("failed to create client directory %s: %w", clientDir, err)
	}

	log.Printf("Generating keys for client %s...", user.TG)
	privateKeyBytes, err := common.RunCommand(ctx, "wg", "genkey")
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("wg genkey for client failed: %w. Stderr: %s", err, string(exitErr.Stderr))
		}
		return fmt.Errorf("wg genkey for client failed: %w", err)
	}
	user.PrivateKey = strings.TrimSpace(string(privateKeyBytes))
	if err := os.WriteFile(clientPrivKeyPath, []byte(user.PrivateKey), 0600); err != nil {
		return fmt.Errorf("failed to write client private key: %w", err)
	}

	cmdPubKey := exec.CommandContext(ctx, "wg", "pubkey")
	cmdPubKey.Stdin = strings.NewReader(user.PrivateKey)
	publicKeyBytes, err := cmdPubKey.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("wg pubkey for client failed: %w. Stderr: %s", err, string(exitErr.Stderr))
		}
		return fmt.Errorf("wg pubkey for client failed: %w", err)
	}
	user.PublicKey = strings.TrimSpace(string(publicKeyBytes))
	if err := os.WriteFile(clientPubKeyPath, []byte(user.PublicKey), 0600); err != nil {
		return fmt.Errorf("failed to write client public key: %w", err)
	}

	clientVPNAddress, err := w.findNextAvailableIP(ctx)
	if err != nil {
		return fmt.Errorf("failed to find available IP for client %s: %w", user.TG, err)
	}
	user.Address = clientVPNAddress
	log.Printf("Assigned IP %s to client %s", user.Address, user.TG)

	log.Printf("Adding client %s to server configuration...", user.TG)

	serverConfBytes, err := os.ReadFile(w.serverConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read server config %s: %w", w.serverConfigPath, err)
	}

	peerComment := fmt.Sprintf("# Client: %s", user.TG)
	if user.TG != "" {
		peerComment += fmt.Sprintf(" (TG: %s)", user.TG)
	}
	if user.Email != "" {
		peerComment += fmt.Sprintf(" (Email: %s)", user.Email)
	}

	peerEntry := fmt.Sprintf("\n[Peer]\n%s\nPublicKey = %s\nAllowedIPs = %s\n",
		peerComment, user.PublicKey, user.Address)

	newServerConf := string(serverConfBytes)
	if !strings.HasSuffix(newServerConf, "\n") && len(newServerConf) > 0 {
		newServerConf += "\n"
	}
	newServerConf += peerEntry

	if err := os.WriteFile(w.serverConfigPath, []byte(newServerConf), 0600); err != nil {
		return fmt.Errorf("failed to append peer to server config %s: %w", w.serverConfigPath, err)
	}

	log.Printf("Creating configuration file for client %s...", user.TG)
	serverPubKeyBytes, err := os.ReadFile(w.serverPubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read server public key for client config: %w", err)
	}
	serverPublicKey := strings.TrimSpace(string(serverPubKeyBytes))

	clientConfData := struct {
		PrivateKey          string
		Address             string
		DNS                 string
		ServerPublicKey     string
		ServerEndpoint      string
		ServerListenPort    int
		PersistentKeepalive int
	}{
		PrivateKey:          user.PrivateKey,
		Address:             user.Address,
		DNS:                 common.DNSAddr,
		ServerPublicKey:     serverPublicKey,
		ServerEndpoint:      w.serverPublicIP,
		ServerListenPort:    w.serverListenPort,
		PersistentKeepalive: 25,
	}

	clientConfTmpl := `[Interface]
PrivateKey = {{.PrivateKey}}
Address = {{.Address}}
DNS = {{.DNS}}

[Peer]
PublicKey = {{.ServerPublicKey}}
Endpoint = {{.ServerEndpoint}}:{{.ServerListenPort}}
AllowedIPs = 0.0.0.0/0,::/0
PersistentKeepalive = {{.PersistentKeepalive}}
` // Added ::/0 for IPv6 if needed, remove if not.
	tmpl, err := template.New("clientConf").Parse(clientConfTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse client config template: %w", err)
	}

	var clientConfBuf bytes.Buffer
	if err := tmpl.Execute(&clientConfBuf, clientConfData); err != nil {
		return fmt.Errorf("failed to execute client config template: %w", err)
	}

	if err := os.WriteFile(clientConfPath, clientConfBuf.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write client config %s: %w", clientConfPath, err)
	}
	log.Printf("Client %s added. Config: %s", user.TG, clientConfPath)

	log.Println("Restarting WireGuard service to apply changes...")
	currentStatus, statusErr := w.GetStatus(ctx)
	if statusErr == nil && currentStatus.IsRunning {
		// wg-quick re-reads config on `wg setconf` which happens internally
		// or a full restart. `wg addconf` or `wg syncconf` is better if available.
		// For wg-quick, often restart is simplest.
		// We can also use `wg syncconf <iface> <(wg-quick strip <iface>)`
		// Or simply reload the service.
		if _, errRel := common.RunCommand(ctx, "systemctl", "reload-or-restart", w.serviceName()); errRel != nil {
			log.Printf("Warning: failed to reload/restart service %s: %v. A manual restart might be needed.", w.serviceName(), errRel)
		} else {
			log.Printf("Service %s reloaded/restarted.", w.serviceName())
		}
	} else if statusErr != nil {
		log.Printf("Warning: could not get service status before restart: %v", statusErr)
	}
	return nil
}

func (w *wireguardInfra) RemoveUser(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	log.Printf("Removing user %s...", userID)

	clientPubKeyPath := filepath.Join(w.clientsDir, userID, "publickey")
	clientPubKeyBytes, err := os.ReadFile(clientPubKeyPath)
	if err != nil {
		return fmt.Errorf("could not read public key for user %s: %w", userID, err)
	}
	clientPublicKey := strings.TrimSpace(string(clientPubKeyBytes))

	log.Printf("Removing peer with PublicKey %s from %s", clientPublicKey, w.serverConfigPath)

	confContentBytes, err := os.ReadFile(w.serverConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read server config %s: %w", w.serverConfigPath, err)
	}

	rePeerSection := regexp.MustCompile(fmt.Sprintf(`(?ms)^\s*\[Peer\][^\[]*?PublicKey\s*=\s*%s.*?(\n\s*\[Peer\]|\n\s*\z|\z)`, regexp.QuoteMeta(clientPublicKey)))
	newConfContent := rePeerSection.ReplaceAllString(string(confContentBytes), "")

	peerFound := newConfContent != string(confContentBytes)
	if !peerFound {
		log.Printf("Warning: Peer section for PublicKey %s not found in %s.", clientPublicKey, w.serverConfigPath)
	} else {
		if err := os.WriteFile(w.serverConfigPath, []byte(strings.TrimSpace(newConfContent)+"\n"), 0600); err != nil {
			return fmt.Errorf("failed to write updated server config %s: %w", w.serverConfigPath, err)
		}
		log.Printf("Peer %s removed from server config.", userID)
	}

	clientDir := filepath.Join(w.clientsDir, userID)
	log.Printf("Removing client directory: %s", clientDir)
	if err := os.RemoveAll(clientDir); err != nil {
		return fmt.Errorf("failed to remove client directory %s: %w", clientDir, err)
	}

	if peerFound {
		log.Println("Reloading/Restarting WireGuard service to apply changes...")
		if _, errRel := common.RunCommand(ctx, "systemctl", "reload-or-restart", w.serviceName()); errRel != nil {
			log.Printf("Warning: failed to reload/restart service %s: %v. A manual restart might be needed.", w.serviceName(), errRel)
		} else {
			log.Printf("Service %s reloaded/restarted.", w.serviceName())
		}
	}
	log.Printf("User %s removed successfully.", userID)
	return nil
}

func (w *wireguardInfra) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User

	confContentBytes, err := os.ReadFile(w.serverConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return users, nil
		}
		return nil, fmt.Errorf("failed to read server config %s: %w", w.serverConfigPath, err)
	}
	content := string(confContentBytes)

	rePeer := regexp.MustCompile(`(?ms)^\s*\[Peer\]\s*\n(?:^\s*#\s*Client:\s*([^\s(]+)(?:\s*\((?:TG:\s*([^)]+))?\s*(?:Email:\s*([^)]+))?\))?.*?\n)?^\s*PublicKey\s*=\s*(\S+)\s*\n^\s*AllowedIPs\s*=\s*(\S+)`)
	matches := rePeer.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		user := &domain.User{}
		if len(match) > 1 && match[1] != "" {
			user.TG = strings.TrimSpace(match[1])
			user.Name = user.TG
		}
		if len(match) > 2 && match[2] != "" {
			user.TG = strings.TrimSpace(match[2])
		}
		if len(match) > 3 && match[3] != "" {
			user.Email = strings.TrimSpace(match[3])
		}
		if len(match) > 4 && match[4] != "" {
			user.PublicKey = strings.TrimSpace(match[4])
		}
		if len(match) > 5 && match[5] != "" {
			user.Address = strings.TrimSpace(match[5])
		}

		if user.PublicKey != "" {
			users = append(users, user)
		}
	}
	return users, nil
}

func (w *wireguardInfra) GetUserTraffic(ctx context.Context, userPublicKey string) (*domain.TrafficInfo, error) {
	if userPublicKey == "" {
		return nil, errors.New("user public key is required")
	}

	dumpOutput, err := common.RunCommand(ctx, "wg", "show", w.wgInterfaceName, "dump")
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'wg show %s dump': %w", w.wgInterfaceName, err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(dumpOutput)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 8 && fields[1] == userPublicKey {
			traffic := &domain.TrafficInfo{}
			if fields[3] != "(none)" {
				host, _, e := net.SplitHostPort(fields[3])
				if e == nil {
					traffic.Endpoint = net.ParseIP(host)
				}
			}
			handshakeSec, e := strconv.ParseInt(fields[5], 10, 64)
			if e == nil && handshakeSec > 0 {
				traffic.LatestHandshake = time.Unix(handshakeSec, 0)
			}
			traffic.ReceivedBytes, _ = strconv.ParseInt(fields[6], 10, 64)
			traffic.TransferBytes, _ = strconv.ParseInt(fields[7], 10, 64)
			return traffic, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning wg dump: %w", err)
	}
	return nil, fmt.Errorf("user %s not found in wg dump", userPublicKey)
}
func (w *wireguardInfra) GetConfig(ctx context.Context, userID string) ([]byte, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	clientConfFileName := userID + ".conf"
	clientConfPath := filepath.Join(w.clientsDir, userID, clientConfFileName)

	configBytes, err := os.ReadFile(clientConfPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fs.ErrNotExist
		}
		return nil, fmt.Errorf("failed to read client config %s: %w", clientConfPath, err)
	}
	return configBytes, nil
}
