package wireguard

import (
	"bufio"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/shanth1/vpn-service/internal/common"
	"github.com/shanth1/vpn-service/internal/core/domain"
)

func (w *wireguardInfra) serviceName() string {
	return fmt.Sprintf("wg-quick@%s", w.wgInterfaceName)
}

func (w *wireguardInfra) StartService(ctx context.Context) error {
	service := w.serviceName()

	_, _ = common.RunCommand(ctx, "systemctl", "enable", service)
	if _, err := common.RunCommand(ctx, "systemctl", "start", service); err != nil {
		return fmt.Errorf("systemctl start: %w", err)
	}

	return nil
}

func (w *wireguardInfra) StopService(ctx context.Context) error {
	service := w.serviceName()

	if _, err := common.RunCommand(ctx, "systemctl", "stop", service); err != nil {
		return fmt.Errorf("systemctl stop: %s", err)
	}

	return nil
}

func (w *wireguardInfra) GetStatus(ctx context.Context) (*domain.Status, error) {
	status := &domain.Status{
		InterfaceName: w.wgInterfaceName,
		ListeningPort: w.serverListenPort,
	}

	service := w.serviceName()
	output, err := common.RunCommand(ctx, "systemctl", "is-active", service)
	if err != nil {
		return nil, fmt.Errorf("systemctl is-active: %w", err)
	}
	status.IsRunning = strings.TrimSpace(string(output)) == "active"

	if status.IsRunning {
		wgShowOutput, err := common.RunCommand(ctx, "wg", "show", w.wgInterfaceName)
		if err != nil {
			return nil, fmt.Errorf("wg show: %w", err)
		}

		scanner := bufio.NewScanner(strings.NewReader(string(wgShowOutput)))
		peerCount := 0
		rePort := regexp.MustCompile(`listening port:\s*(\d+)`)
		rePeer := regexp.MustCompile(`peer:\s*(\S+)`)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if matches := rePort.FindStringSubmatch(line); len(matches) > 1 {
				if port, errConv := strconv.Atoi(matches[1]); errConv == nil {
					status.ListeningPort = port
				}
			}
			if rePeer.MatchString(line) {
				peerCount++
			}
		}
		status.Peers = peerCount
	}

	return status, nil
}
