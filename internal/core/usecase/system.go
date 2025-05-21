package usecase

import (
	"context"
)

func (s *vpnCoreImpl) Install(ctx context.Context) error {
	return s.protocol.Install(ctx)
}

func (s *vpnCoreImpl) Uninstall(ctx context.Context) error {
	return s.protocol.Uninstall(ctx)
}
