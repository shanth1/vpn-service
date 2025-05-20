package system

import "context"

func (s *fakeSystemInfra) GetNetworkInfo() (ifaceName, ip string, err error) {
	return "eth0", "10.0.0.1", nil
}

func (s *fakeSystemInfra) SetUpRedirection(ctx context.Context) error {
	return nil
}
