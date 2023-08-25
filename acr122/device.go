package acr122

import "github.com/maitredede/gonfc"

type LedState struct {
	Red   bool
	Green bool
}

type ACR122Device interface {
	gonfc.Device

	Name() string
	GetLedState() (LedState, error)
	SetLed(state LedState) (LedState, error)
}
