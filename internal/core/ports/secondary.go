package ports

import (
	"context"

	"github.com/shanth1/vpn-service/internal/core/domain"
)

// Protocol is secondary port for working with vpn tool
type Protocol interface {
	Install(ctx context.Context) error
	Uninstall(ctx context.Context) error

	SetUpServer(ctx context.Context) error
	TearDownServer(ctx context.Context) error

	StartService(ctx context.Context) error
	StopService(ctx context.Context) error
	GetStatus(ctx context.Context) (*domain.Status, error)

	AddUser(ctx context.Context, user *domain.User) error
	RemoveUser(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserTraffic(ctx context.Context, userPublicKey string) (*domain.TrafficInfo, error)
	GetConfig(ctx context.Context, id string) ([]byte, error)
}

// QRCode is secondary port for working with qr codes
type QRCode interface {
	Generate(ctx context.Context, outputDir, fileName string, data []byte) error
}

// System is secondary port for working with system
type System interface {
	SetUpRedirection(ctx context.Context) error
	GetNetworkInfo() (ifaceName, ip string, err error)
}
