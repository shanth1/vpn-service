package system

import "context"

func (s *fakeSystemAdapter) GetNetworkInfo() (ifaceName, ip string, err error) {
	return "eth0", "10.0.0.1", nil
}

func (s *fakeSystemAdapter) SetUpRedirection(ctx context.Context) error {
	return nil
}
