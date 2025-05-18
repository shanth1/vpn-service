package wireguard

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/shanth1/vpn-service/internal/core/domain"
)

func (w *wireguardInfra) findNextAvailableIP(ctx context.Context) (string, error) {
	_, ipNet, err := net.ParseCIDR(w.serverVPNAddrCIDR)
	if err != nil {
		return "", fmt.Errorf("invalid server VPN address format '%s': %w", w.serverVPNAddrCIDR, err)
	}

	existingUsers, err := w.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Warning: could not get all users to determine next IP: %v. Will try to find first available.", err)
		existingUsers = []*domain.User{}
	}

	existingIPs := make(map[string]bool)
	serverIPOnly := strings.Split(w.serverVPNAddrCIDR, "/")[0]
	existingIPs[serverIPOnly] = true

	for _, user := range existingUsers {
		if user.Address != "" {
			userIPOnly := strings.Split(user.Address, "/")[0]
			existingIPs[userIPOnly] = true
		}
	}

	currentIP := make(net.IP, len(ipNet.IP))
	copy(currentIP, ipNet.IP)

	// Начинаем поиск со следующего IP после серверного (или с .2, если сервер .1)
	// Если сервер, например, 10.8.0.1, начинаем с 10.8.0.2
	// Для этого увеличиваем IP на 1 до тех пор, пока он не станет больше серверного IP (если сервер не .0)
	serverNumIP := net.ParseIP(serverIPOnly)
	if serverNumIP == nil {
		return "", fmt.Errorf("could not parse server IP %s", serverIPOnly)
	}

	// Инкремент currentIP, пока оно не будет > serverNumIP или пока мы не переберем всю сеть
	for i := 0; i < 255; i++ { // Ограничение, чтобы избежать бесконечного цикла на маленьких сетях
		// Увеличиваем IP на 1
		for j := len(currentIP) - 1; j >= 0; j-- {
			currentIP[j]++
			if currentIP[j] > 0 {
				break
			}
		}
		if !ipNet.Contains(currentIP) { // Вышли за пределы подсети
			return "", errors.New("no available IP addresses in the subnet (exhausted range)")
		}
		if bytes.Compare(currentIP, serverNumIP) <= 0 { // Пропускаем IP сервера и меньшие
			continue
		}

		// Проверка на broadcast (X.Y.Z.255 для /24)
		ones, bits := ipNet.Mask.Size()
		if ones == bits-8 && currentIP[len(currentIP)-1] == 255 { // Простая проверка для /24, /16, /8
			continue
		}

		candidateIP := currentIP.String()
		if !existingIPs[candidateIP] {
			return fmt.Sprintf("%s/32", candidateIP), nil
		}
	}
	return "", errors.New("no available IP addresses in the subnet (all occupied or too small range)")
}
