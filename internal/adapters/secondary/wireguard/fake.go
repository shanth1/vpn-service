package wireguard

import (
	"context"
	"errors"
	"time"

	"github.com/shanth1/vpn-service/internal/core/domain"
)

var (
	userDB = map[string]*domain.User{
		"user": {
			TG:         "user",
			PublicKey:  "FdjOER23",
			PrivateKey: "jdowdfefFFD",
			Name:       "Ivan",
			Surname:    "Ivanov",
			Email:      "ivan@mail.ru",
			Address:    "172.172.172.172",
		},
		"elpepa": {
			TG:         "elpepa",
			PublicKey:  "jfFDkFO",
			PrivateKey: "NVCIEee",
			Name:       "Igor",
			Address:    "100.200.101.1",
		},
	}
)

const (
	fakeExecutionTime = 2 * time.Second
)

func (a *wgFakeAdapter) SetUpServer(ctx context.Context) error {
	time.Sleep(fakeExecutionTime)
	return nil
}
func (a *wgFakeAdapter) TearDownServer(ctx context.Context) error {
	time.Sleep(fakeExecutionTime)
	return nil
}

func (a *wgFakeAdapter) StartService(ctx context.Context) error {
	time.Sleep(fakeExecutionTime)
	return nil
}

func (a *wgFakeAdapter) StopService(ctx context.Context) error {
	time.Sleep(fakeExecutionTime)
	return nil
}

func (a *wgFakeAdapter) GetStatus(ctx context.Context) (*domain.Status, error) {
	time.Sleep(fakeExecutionTime)
	return &domain.Status{}, nil
}

func (*wgFakeAdapter) Install(ctx context.Context) error {
	time.Sleep(fakeExecutionTime)
	return nil
}

func (*wgFakeAdapter) Uninstall(ctx context.Context) error {
	time.Sleep(fakeExecutionTime)
	return nil
}

func (a *wgFakeAdapter) AddUser(ctx context.Context, user *domain.User) error {
	time.Sleep(fakeExecutionTime)
	if _, ok := userDB[user.TG]; ok {
		return errors.New("already exists")
	}
	userDB[user.TG] = user
	return nil
}

func (a *wgFakeAdapter) RemoveUser(ctx context.Context, userID string) error {
	time.Sleep(fakeExecutionTime)
	if _, ok := userDB[userID]; !ok {
		return errors.New("user not found")
	}
	delete(userDB, userID)
	return nil
}

func (a *wgFakeAdapter) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	time.Sleep(fakeExecutionTime)
	var users []*domain.User
	for _, user := range userDB {
		users = append(users, user)
	}
	return users, nil
}

func (a *wgFakeAdapter) GetUserTraffic(ctx context.Context, userPublicKey string) (*domain.TrafficInfo, error) {
	return nil, errors.New("not implemented")
}

func (a *wgFakeAdapter) GetConfig(ctx context.Context, userID string) ([]byte, error) {
	return nil, errors.New("not implemented")
}
