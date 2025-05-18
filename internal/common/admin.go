package common

import (
	"os"
)

func CheckAdminRights() bool {
	if os.Geteuid() != 0 {
		return false
	}

	return true
}
