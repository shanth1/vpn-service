package wireguard

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/shanth1/vpn-service/internal/common"
)

// TODO: move to system adapter
func (w *wireguardInfra) SetUpRedirection(ctx context.Context) error {
	if w.publicNetIface == "" {
		return errors.New("empy net interface")
	}

	log.Println("Enabling IPv4 forwarding...")
	if _, err := common.RunCommand(ctx, "sysctl", "-w", "net.ipv4.ip_forward=1"); err != nil {
		return fmt.Errorf("enable ipv4 forwarding: %w", err)
	}

	sysctlConfPath := "/etc/sysctl.conf"
	sysctlLine := "net.ipv4.ip_forward=1"

	contentBytes, err := os.ReadFile(sysctlConfPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read sysctl config: %w", sysctlConfPath, err)
	}
	content := string(contentBytes)
	re := regexp.MustCompile(`(?m)^\s*#?\s*net.ipv4.ip_forward\s*=\s*.*`)
	if !re.MatchString(content) {
		if content != "" && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += sysctlLine + "\n"
	} else {
		content = re.ReplaceAllString(content, sysctlLine)
	}
	if err := os.WriteFile(sysctlConfPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write updated content: %w", sysctlConfPath, err)
	}
	if _, err := common.RunCommand(ctx, "sysctl", "-p"); err != nil {
		return fmt.Errorf("sysctl update: %w", err)
	}

	return nil
}
