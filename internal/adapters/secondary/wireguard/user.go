package wireguard

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/shanth1/vpn-service/internal/common"
	"github.com/shanth1/vpn-service/internal/core/domain"
)

func (a *wgAdapter) AddUser(ctx context.Context, user *domain.User) error {
	if user.TG == "" {
		return errors.New("empty tg field")
	}

	if a.serverPublicIP == "" {
		return errors.New("empy public ip")
	}

	clientDir := filepath.Join(a.clientsDir, user.TG)

	if err := os.MkdirAll(clientDir, 0700); err != nil {
		return fmt.Errorf("create client directory: %w", err)
	}

	privateKeyBytes, publicKeyBytes, err := genWgKeys(ctx)
	if err != nil {
		return fmt.Errorf("gen keys: %w", err)
	}

	if err := os.WriteFile(getPrivateKeyPath(clientDir), []byte(user.PrivateKey), 0600); err != nil {
		return fmt.Errorf("write client private key: %w", err)
	}
	user.PrivateKey = string(privateKeyBytes)

	if err := os.WriteFile(getPublicKeyPath(clientDir), []byte(user.PublicKey), 0600); err != nil {
		return fmt.Errorf("write client public key: %w", err)
	}
	user.PublicKey = strings.TrimSpace(string(publicKeyBytes))

	clientVPNAddress, err := a.findNextAvailableIP(ctx)
	if err != nil {
		return fmt.Errorf("find available IP: %w", err)
	}
	user.Address = clientVPNAddress

	serverConfBytes, err := os.ReadFile(a.serverConfigPath)
	if err != nil {
		return fmt.Errorf("read server config: %w", err)
	}

	peerComment := fmt.Sprintf("# Client: %s", user.TG)
	if user.TG != "" {
		peerComment += fmt.Sprintf(" (TG: %s)", user.TG)
	}
	if user.Name != "" {
		peerComment += fmt.Sprintf(" (Name: %s)", user.TG)
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

	if err := os.WriteFile(a.serverConfigPath, []byte(newServerConf), 0600); err != nil {
		return fmt.Errorf("append peer to server config: %w", err)
	}

	serverPubKeyBytes, err := os.ReadFile(a.serverPubKeyPath)
	if err != nil {
		return fmt.Errorf("read server public key: %w", err)
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
		ServerEndpoint:      a.serverPublicIP,
		ServerListenPort:    a.serverListenPort,
		PersistentKeepalive: 25,
	}

	clientConfTmpl := `[Interface]
PrivateKey = {{.PrivateKey}}
Address = {{.Address}}
DNS = {{.DNS}}

[Peer]
PublicKey = {{.ServerPublicKey}}
Endpoint = {{.ServerEndpoint}}:{{.ServerListenPort}}
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = {{.PersistentKeepalive}}
`
	tmpl, err := template.New("clientConf").Parse(clientConfTmpl)
	if err != nil {
		return fmt.Errorf("parse config template: %w", err)
	}

	var clientConfBuf bytes.Buffer
	if err := tmpl.Execute(&clientConfBuf, clientConfData); err != nil {
		return fmt.Errorf("execute config template: %w", err)
	}

	clientConfFilePath := getConfFileName("client")
	if err := os.WriteFile(clientConfFilePath, clientConfBuf.Bytes(), 0600); err != nil {
		return fmt.Errorf("write client config: %w", err)
	}

	currentStatus, _ := a.GetStatus(ctx)
	if err != nil {
		return fmt.Errorf("get status: %w", err)
	}

	if currentStatus.IsRunning {
		// wg-quick re-reads config on `wg setconf` which happens internally
		// or a full restart. `wg addconf` or `wg syncconf` is better if available.
		// For wg-quick, often restart is simplest.
		// We can also use `wg syncconf <iface> <(wg-quick strip <iface>)`
		// Or simply reload the service.
		if _, errRel := common.RunCommand(ctx, "systemctl", "reload-or-restart", a.serviceName()); errRel != nil {
			return fmt.Errorf("reload or restart: %w", err)
		}
	}

	return nil
}

func (a *wgAdapter) RemoveUser(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	clientDir := filepath.Join(a.clientsDir, userID)

	clientPublicKey, err := getKey(getPublicKeyPath(clientDir))
	if err != nil {
		return fmt.Errorf("get public key: %w", userID, err)
	}

	confContentBytes, err := os.ReadFile(a.serverConfigPath)
	if err != nil {
		return fmt.Errorf("read server config %s: %w", a.serverConfigPath, err)
	}

	rePeerSection := regexp.MustCompile(fmt.Sprintf(`(?ms)^\s*\[Peer\][^\[]*?PublicKey\s*=\s*%s.*?(\n\s*\[Peer\]|\n\s*\z|\z)`, regexp.QuoteMeta(clientPublicKey)))
	newConfContent := rePeerSection.ReplaceAllString(string(confContentBytes), "")

	peerFound := newConfContent != string(confContentBytes)
	if !peerFound {
		return fmt.Errorf("peer not found: %w", err)
	} else {
		if err := os.WriteFile(a.serverConfigPath, []byte(strings.TrimSpace(newConfContent)+"\n"), 0600); err != nil {
			return fmt.Errorf("failed to write updated server config %s: %w", a.serverConfigPath, err)
		}
	}

	if err := os.RemoveAll(clientDir); err != nil {
		return fmt.Errorf("remove client directory: %w", err)
	}

	if peerFound {
		if _, errRel := common.RunCommand(ctx, "systemctl", "reload-or-restart", a.serviceName()); errRel != nil {
			return fmt.Errorf("reload or restart: %w", err)
		}
	}
	return nil
}

func (a *wgAdapter) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User

	confContentBytes, err := os.ReadFile(a.serverConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return users, nil
		}
		return nil, fmt.Errorf("read server config: %w", err)
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

func (a *wgAdapter) GetUserTraffic(ctx context.Context, userPublicKey string) (*domain.TrafficInfo, error) {
	if userPublicKey == "" {
		return nil, errors.New("empy public key")
	}

	dumpOutput, err := common.RunCommand(ctx, "wg", "show", a.wgInterfaceName, "dump")
	if err != nil {
		return nil, fmt.Errorf("wg show: %w", err)
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
		return nil, fmt.Errorf("scanning wg dump: %w", err)
	}
	return nil, fmt.Errorf("user not found in wg dump")
}

func (a *wgAdapter) GetConfig(ctx context.Context, userID string) ([]byte, error) {
	if userID == "" {
		return nil, errors.New("empty user id")
	}

	clientConfPath := filepath.Join(a.clientsDir, userID, getConfFileName("client"))

	configBytes, err := os.ReadFile(clientConfPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fs.ErrNotExist
		}
		return nil, fmt.Errorf("read client config: %w", err)
	}
	return configBytes, nil
}
