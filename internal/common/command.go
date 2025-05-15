package common

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(ctx context.Context, name string, arg ...string) error {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("command %s %v cancelled: %w", name, arg, ctx.Err())
		}
		return fmt.Errorf("command %s %v failed: %w", name, arg, err)
	}

	return nil
}
