package system

type systemAdapter struct{}

func NewAdapter() *systemAdapter {
	return &systemAdapter{}
}

type fakeSystemAdapter struct{}

func NewFakeAdapter() *fakeSystemAdapter {
	return &fakeSystemAdapter{}
}
