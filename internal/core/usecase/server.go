package usecase

import (
	"context"
	"fmt"
)

func (s *vpnServiceImpl) Init(ctx context.Context) error {
	if err := s.protocol.SetUpServer(ctx); err != nil {
		return fmt.Errorf("server: %w", err)
	}

	if err := s.protocol.SetUpRedirection(ctx); err != nil {
		return fmt.Errorf("redirection: %w", err)
	}

	return nil
}

func (s *vpnServiceImpl) CleanUp(ctx context.Context) error {
	return s.protocol.TearDownServer(ctx)
}
