package wireguard

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/shanth1/vpn-service/internal/core/domain"
)

func (a *wgAdapter) findNextAvailableIP(ctx context.Context) (string, error) {
	_, ipNet, err := net.ParseCIDR(a.serverVPNAddrCIDR)
	if err != nil {
		return "", fmt.Errorf("invalid server VPN address format '%s': %w", a.serverVPNAddrCIDR, err)
	}

	existingUsers, err := a.GetAllUsers(ctx)
	if err != nil {
		existingUsers = []*domain.User{}
	}

	existingIPs := make(map[string]bool)
	serverIPOnly := strings.Split(a.serverVPNAddrCIDR, "/")[0]
	existingIPs[serverIPOnly] = true

	for _, user := range existingUsers {
		if user.Address != "" {
			userIPOnly := strings.Split(user.Address, "/")[0]
			existingIPs[userIPOnly] = true
		}
	}

	currentIP := make(net.IP, len(ipNet.IP))
	copy(currentIP, ipNet.IP)

	serverNumIP := net.ParseIP(serverIPOnly)
	if serverNumIP == nil {
		return "", fmt.Errorf("parse server IP: %s", serverIPOnly)
	}

	for i := 0; i < 255; i++ {
		for j := len(currentIP) - 1; j >= 0; j-- {
			currentIP[j]++
			if currentIP[j] > 0 {
				break
			}
		}
		if !ipNet.Contains(currentIP) {
			return "", errors.New("no available IP addresses in the subnet (exhausted range)")
		}
		if bytes.Compare(currentIP, serverNumIP) <= 0 {
			continue
		}

		ones, bits := ipNet.Mask.Size()
		if ones == bits-8 && currentIP[len(currentIP)-1] == 255 {
			continue
		}

		candidateIP := currentIP.String()
		if !existingIPs[candidateIP] {
			return fmt.Sprintf("%s/32", candidateIP), nil
		}
	}
	return "", errors.New("no available IP addresses in the subnet (all occupied or too small range)")
}
