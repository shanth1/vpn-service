package ports

import (
	"context"

	"github.com/shanth1/vpn-service/internal/core/domain"
)

type PrimaryPort interface {
	// System group

	Install(ctx context.Context) error   // Install installs all necessary packages
	Uninstall(ctx context.Context) error // Uninstall removes all used packages
	// ---------------------

	// Server group

	Init(ctx context.Context) error    // Init creates and configures a server
	CleanUp(ctx context.Context) error // CleanUp deletes server files
	// ---------------------

	// Service

	StartService(ctx context.Context) error
	StopService(ctx context.Context) error
	GetStatus(ctx context.Context) (*domain.Status, error)
	// ---------------------

	// User group

	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	AddUser(ctx context.Context, user *domain.User) error
	RemoveUser(ctx context.Context, id string) error
	SaveConfig(ctx context.Context, id string) error // GetConfig Saves the config file and QR code in place of the binary
}
