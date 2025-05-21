package usecase

import (
	"context"
	"fmt"

	"github.com/shanth1/vpn-service/internal/core/domain"
)

func (s *vpnCoreImpl) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.protocol.GetAllUsers(ctx)
}

func (s *vpnCoreImpl) AddUser(ctx context.Context, user *domain.User) error {
	return s.protocol.AddUser(ctx, user)
}

func (s *vpnCoreImpl) RemoveUser(ctx context.Context, id string) error {
	return s.protocol.RemoveUser(ctx, id)
}

func (s *vpnCoreImpl) SaveConfig(ctx context.Context, id string) error {
	config, err := s.protocol.GetConfig(ctx, id)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// TODO: save config

	if err := s.qr.Generate(ctx, ".", id, config); err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	return nil
}
