package wireguard

import (
	"context"
	"fmt"

	"github.com/shanth1/vpn-service/internal/common"
)

func (*wgAdapter) Install(ctx context.Context) error {
	if _, err := common.RunCommand(ctx, "apt", "update", "-y"); err != nil {
		return fmt.Errorf("update packages: %w", err)
	}

	_, _ = common.RunCommand(ctx, "apt", "upgrade", "-y")

	if _, err := common.RunCommand(ctx, "apt", "install", "-y", "wireguard", "wireguard-tools"); err != nil {
		return fmt.Errorf("install packages: %w", err)
	}

	return nil
}

func (*wgAdapter) Uninstall(ctx context.Context) error {
	if _, err := common.RunCommand(ctx, "apt", "remove", "-y", "wireguard", "wireguard-tools"); err != nil {
		return fmt.Errorf("remove packages: %w", err)
	}

	_, _ = common.RunCommand(ctx, "apt", "autoremove", "-y")

	return nil
}
