package wireguard

import (
	"bufio"
	"context"
	"fmt"
	"log"
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
	log.Printf("Enabling %s service...", service)
	_, _ = common.RunCommand(ctx, "systemctl", "enable", service)

	log.Printf("Starting %s service...", service)
	if _, err := common.RunCommand(ctx, "systemctl", "start", service); err != nil {
		return fmt.Errorf("systemctl start: %w", err)
	}
	log.Printf("%s service started.", service)
	return nil
}

func (w *wireguardInfra) StopService(ctx context.Context) error {
	service := w.serviceName()
	log.Printf("Stopping %s service...", service)
	if _, err := common.RunCommand(ctx, "systemctl", "stop", service); err != nil {
		return fmt.Errorf("systemctl stop: %s", err)
	}
	log.Printf("%s service stopped.", service)
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
		if strings.Contains(err.Error(), "exit status") {
			sOutput := strings.TrimSpace(string(output))
			if sOutput == "inactive" || sOutput == "failed" {
				status.IsRunning = false
			} else {
				log.Printf("%s service status check (is-active) failed: %v, output: %s", service, err, sOutput)
				status.IsRunning = false
			}
		} else {
			log.Printf("Error executing systemctl is-active for %s: %v", service, err)
			status.IsRunning = false
		}
	} else {
		status.IsRunning = strings.TrimSpace(string(output)) == "active"
	}

	if !status.IsRunning {
		log.Printf("%s service is not running, skipping wg show.", service)
		return status, nil
	}

	wgShowOutput, err := common.RunCommand(ctx, "wg", "show", w.wgInterfaceName)
	if err != nil {
		log.Printf("Failed to execute 'wg show %s': %v. Status might be incomplete.", w.wgInterfaceName, err)
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

	return status, nil
}
