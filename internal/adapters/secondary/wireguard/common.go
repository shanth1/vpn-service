package wireguard

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shanth1/vpn-service/internal/common"
)

// getConfFileName returns filename with .conf
func getConfFileName(fileName string) string {
	return fileName + ".conf"
}

func getPrivateKeyPath(dir string) string {
	return filepath.Join(dir, "privatekey")
}

func getPublicKeyPath(dir string) string {
	return filepath.Join(dir, "publickey")
}

func getKey(path string) (string, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(keyBytes)), nil
}

func genWgKeys(ctx context.Context) (privateKeyBytes, publicKeyBytes []byte, err error) {
	privateKeyBytes, err = common.RunCommand(ctx, "wg", "genkey")
	if err != nil {
		return nil, nil, fmt.Errorf("wg genkey: %w", err)
	}
	privateKeyBytes = bytes.TrimSpace(privateKeyBytes)

	cmdPubKey := exec.CommandContext(ctx, "wg", "pubkey")
	cmdPubKey.Stdin = bytes.NewReader(privateKeyBytes)
	publicKeyBytes, err = cmdPubKey.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, nil, fmt.Errorf("wg pubkey: %w. Stderr: %s", err, string(exitErr.Stderr))
		}
		return nil, nil, fmt.Errorf("wg pubkey: %w", err)
	}
	publicKeyBytes = bytes.TrimSpace(publicKeyBytes)

	return privateKeyBytes, publicKeyBytes, nil
}
