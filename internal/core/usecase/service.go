package usecase

import (
	"context"

	"github.com/shanth1/vpn-service/internal/core/domain"
)

func (s *vpnCoreImpl) StartService(ctx context.Context) error {
	return s.protocol.StartService(ctx)
}

func (s *vpnCoreImpl) StopService(ctx context.Context) error {
	return s.protocol.StopService(ctx)
}

func (s *vpnCoreImpl) GetStatus(ctx context.Context) (*domain.Status, error) {
	return s.protocol.GetStatus(ctx)
}
