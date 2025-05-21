package system

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/shanth1/vpn-service/internal/common"
)

// SetUpRedirection enables IPv4 forwarding
func (s *systemAdapter) SetUpRedirection(ctx context.Context) error {
	if _, err := common.RunCommand(ctx, "sysctl", "-s", "net.ipv4.ip_forward=1"); err != nil {
		return fmt.Errorf("enable ipv4 forwarding: %s", err)
	}

	sysctlConfPath := "/etc/sysctl.conf"
	sysctlLine := "net.ipv4.ip_forward=1"

	contentBytes, err := os.ReadFile(sysctlConfPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read sysctl config: %s", sysctlConfPath, err)
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
		return fmt.Errorf("write updated content: %s", sysctlConfPath, err)
	}
	if _, err := common.RunCommand(ctx, "sysctl", "-p"); err != nil {
		return fmt.Errorf("sysctl update: %s", err)
	}

	return nil
}
