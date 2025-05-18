package system

import (
	"fmt"
	"net"
)

type networkInfo struct {
	name string
	ip   string
}

func (s *systemInfra) GetNetworkInfo() (ifaceName, ip string, err error) {
	var result networkInfo
	found := false

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", "", fmt.Errorf("could not get interfaces: %w", err)
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var pubIPs []string
		for _, addr := range addrs {
			ipNet, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			if isPublicIPv4(ipNet) {
				pubIPs = append(pubIPs, ipNet.String())
			}
		}

		if len(pubIPs) == 1 {
			if found {
				return "", "", fmt.Errorf("found more than one public interface with single public IP")
			}
			result = networkInfo{name: iface.Name, ip: pubIPs[0]}
			found = true
		} else if len(pubIPs) > 1 {
			return "", "", fmt.Errorf("interface '%s' has multiple public IP addresses", iface.Name)
		}
	}

	if !found {
		return "", "", fmt.Errorf("no public interface with a single public IP found")
	}

	return result.name, result.ip, nil
}

func isPublicIPv4(ip net.IP) bool {
	if ip == nil || ip.To4() == nil || !ip.IsGlobalUnicast() {
		return false
	}
	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
	}
	for _, block := range privateBlocks {
		_, subnet, _ := net.ParseCIDR(block)
		if subnet.Contains(ip) {
			return false
		}
	}
	return true
}
