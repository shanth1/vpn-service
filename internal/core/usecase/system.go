package usecase

import (
	"context"
)

func (s *vpnServiceImpl) Install(ctx context.Context) error {
	return s.protocol.Install(ctx)
}

func (s *vpnServiceImpl) Uninstall(ctx context.Context) error {
	return s.protocol.Uninstall(ctx)
}
