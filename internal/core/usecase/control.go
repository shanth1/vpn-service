package usecase

import (
	"context"

	"github.com/shanth1/vpn-service/internal/core/domain"
	"github.com/shanth1/vpn-service/internal/core/ports"
)

type vpnServiceImpl struct {
	// protocol ports.Protocol
}

func NewVPNService() ports.PrimaryPort {
	return &vpnServiceImpl{
		// protocol: protocol,
	}
}

func (s *vpnServiceImpl) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return []*domain.User{}, nil
}
