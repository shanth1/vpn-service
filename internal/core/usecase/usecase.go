package usecase

import (
	"github.com/shanth1/vpn-service/internal/core/ports"
)

type vpnServiceImpl struct {
	protocol ports.Protocol
	qr       ports.QRCode
}

func NewVPNService(protocolInfra ports.Protocol, qrInfra ports.QRCode) ports.PrimaryPort {
	return &vpnServiceImpl{
		protocol: protocolInfra,
		qr:       qrInfra,
	}
}
