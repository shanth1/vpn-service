package usecase

import (
	"github.com/shanth1/vpn-service/internal/core/ports"
)

type vpnServiceImpl struct {
	protocol ports.Protocol
	qr       ports.QRCode
}

func NewVPNService(protocol ports.Protocol, qr ports.QRCode) ports.PrimaryPort {
	return &vpnServiceImpl{
		protocol: protocol,
		qr:       qr,
	}
}
