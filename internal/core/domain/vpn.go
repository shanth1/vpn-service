package domain

type User struct {
	ID   string
	Name string
}

type TrafficInfo struct {
	UserID   string
	Upload   int64
	Download int64
	LastSeen string
}

type Status struct {
	IsRunning  bool
	Peers      int
	ListenPort int
}
