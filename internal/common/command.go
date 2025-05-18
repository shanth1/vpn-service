package common

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func RunCommand(ctx context.Context, name string, arg ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, arg...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.Canceled {
			return stdout.Bytes(), fmt.Errorf("command %q %v cancelled: %w", name, arg, ctx.Err())
		}
		return stdout.Bytes(), fmt.Errorf("command %q %v failed: %w. Stderr: %s", name, arg, err, stderr.String())
	}
	return stdout.Bytes(), nil
}
