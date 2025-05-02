package ports

import (
	"context"

	"github.com/shanth1/vpn-service/internal/core/domain"
)

type PrimaryPort interface {
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	// GetUserTraffic(ctx context.Context, userID string) (domain.TrafficInfo, error)
	// GetAllTrafficInfo(ctx context.Context) (map[string]domain.TrafficInfo, error)

	// GetStatus(ctx context.Context) (domain.Status, error)

	// Init(ctx context.Context) error
	// RegisterService(ctx context.Context) error
	// DeleteService(ctx context.Context) error
	// StartService(ctx context.Context) error
	// RestartService(ctx context.Context) error
	// StopService(ctx context.Context) error

	// AddUser(ctx context.Context, username string) error
	// UpdateUser(ctx context.Context, username string) error
	// RemoveUser(ctx context.Context, username string) error
}

type Protocol interface {
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	// GetUserTraffic(ctx context.Context, userID string) (domain.TrafficInfo, error)
	// GetAllTrafficInfo(ctx context.Context) (map[string]domain.TrafficInfo, error)

	// GetStatus(ctx context.Context) (domain.Status, error)

	// Init(ctx context.Context) error
	// RegisterService(ctx context.Context) error
	// DeleteService(ctx context.Context) error
	// StartService(ctx context.Context) error
	// StopService(ctx context.Context) error

	// AddUser(ctx context.Context, username string) error
	// RemoveUser(ctx context.Context, username string) error
}
