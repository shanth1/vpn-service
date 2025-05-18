package wireguard

import (
	"context"
	"fmt"

	"github.com/shanth1/vpn-service/internal/common"
)

func (w *wireguardInfra) Install(ctx context.Context) error {
	if _, err := common.RunCommand(ctx, "apt", "update", "-y"); err != nil {
		return fmt.Errorf("update packages: %w", err)
	}

	if _, err := common.RunCommand(ctx, "apt", "upgrade", "-y"); err != nil {
		// TODO: log warn
	}

	if _, err := common.RunCommand(ctx, "apt", "install", "-y", "wireguard", "wireguard-tools"); err != nil {
		return fmt.Errorf("install wireguard packages: %w", err)
	}

	return nil
}

func (w *wireguardInfra) Uninstall(ctx context.Context) error {
	if _, err := common.RunCommand(ctx, "apt", "remove", "-y", "wireguard", "wireguard-tools"); err != nil {
		return fmt.Errorf("failed to remove wireguard packages: %w", err)
	}

	if _, err := common.RunCommand(ctx, "apt", "autoremove", "-y"); err != nil {
		// TODO: log warn
	}

	return nil
}
