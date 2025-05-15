package usecase

import (
	"context"
	"fmt"
)

func (s *vpnServiceImpl) Install(ctx context.Context) error {
	if err := s.protocol.Install(ctx); err != nil {
		return fmt.Errorf("protocol: %w", err)
	}

	if err := s.qr.Install(ctx); err != nil {
		return fmt.Errorf("qr: %w", err)
	}

	return nil
}

func (s *vpnServiceImpl) Uninstall(ctx context.Context) error {
	if err := s.protocol.Uninstall(ctx); err != nil {
		return fmt.Errorf("protocol: %w", err)
	}

	if err := s.qr.Uninstall(ctx); err != nil {
		return fmt.Errorf("qr: %w", err)
	}

	return nil
}
