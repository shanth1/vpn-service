package users

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shanth1/vpn-service/internal/core/domain"
	"github.com/shanth1/vpn-service/internal/core/ports"
)

type allUsersLoadedMsg struct {
	users []*domain.User
	err   error
}
type userDetailsLoadedMsg struct {
	traffic *domain.TrafficInfo
	err     error
}
type userDeletedMsg struct {
	userID string
	err    error
}
type userCreatedMsg struct {
	user *domain.User
	err  error
}

func fetchAllUsersCmd(core ports.PrimaryPort) tea.Cmd {
	return func() tea.Msg {
		users, err := core.GetAllUsers(context.Background())
		return allUsersLoadedMsg{users: users, err: err}
	}
}

func fetchUserDetailsCmd(core ports.PrimaryPort, userPublicKey string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second)

		fmt.Println(core, userPublicKey) // TODO: remove

		// TODO: added method for fetching user traffic
		traffic := &domain.TrafficInfo{}
		return userDetailsLoadedMsg{traffic: traffic, err: nil}
	}
}

func deleteUserCmd(core ports.PrimaryPort, userID string) tea.Cmd {
	return func() tea.Msg {
		err := core.RemoveUser(context.Background(), userID)
		return userDeletedMsg{userID: userID, err: err}
	}
}

func createUserCmd(core ports.PrimaryPort, user *domain.User) tea.Cmd {
	return func() tea.Msg {
		err := core.AddUser(context.Background(), user)
		return userCreatedMsg{user: user, err: err}
	}
}
