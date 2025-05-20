package system

type systemInfra struct{}

func NewInfra() *systemInfra {
	return &systemInfra{}
}

type fakeSystemInfra struct{}

func NewFakeInfra() *fakeSystemInfra {
	return &fakeSystemInfra{}
}
