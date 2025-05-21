package usecase

import (
	"github.com/shanth1/vpn-service/internal/core/ports"
)

type vpnCoreImpl struct {
	protocol ports.Protocol
	qr       ports.QRCode
	system   ports.System
}

func NewVPNCore(
	systemAdapter ports.System,
	protocolAdapter ports.Protocol,
	qrAdapter ports.QRCode,
) ports.PrimaryPort {
	return &vpnCoreImpl{
		system:   systemAdapter,
		protocol: protocolAdapter,
		qr:       qrAdapter,
	}
}
