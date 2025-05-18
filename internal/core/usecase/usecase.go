package usecase

import (
	"github.com/shanth1/vpn-service/internal/core/ports"
)

type vpnServiceImpl struct {
	protocol ports.Protocol
	qr       ports.QRCode
	system   ports.System
}

func NewVPNService(
	systemInfra ports.System,
	protocolInfra ports.Protocol,
	qrInfra ports.QRCode,
) ports.PrimaryPort {
	return &vpnServiceImpl{
		system:   systemInfra,
		protocol: protocolInfra,
		qr:       qrInfra,
	}
}
