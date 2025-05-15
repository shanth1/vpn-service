package domain

import (
	"net"
	"time"
)

type User struct {
	TG        string
	PublicKey string
	Name      string
	Surname   string
	Email     string
}

type TrafficInfo struct {
	Endpoint        net.IP // user ip
	TransferBytes   int64
	ReceivedBytes   int64
	LatestHandshake time.Time
}

type Status struct {
	IsRunning     bool
	Peers         int
	ListeningPort int
	Interface     *net.Interface
}
