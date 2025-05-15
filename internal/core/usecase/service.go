package usecase

import (
	"context"

	"github.com/shanth1/vpn-service/internal/core/domain"
)

func (s *vpnServiceImpl) StartService(ctx context.Context) error {
	return s.protocol.StartService(ctx)
}

func (s *vpnServiceImpl) StopService(ctx context.Context) error {
	return s.protocol.StopService(ctx)
}

func (s *vpnServiceImpl) GetStatus(ctx context.Context) (*domain.Status, error) {
	return s.protocol.GetStatus(ctx)
}
