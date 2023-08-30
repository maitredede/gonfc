package gonfc

import (
	"time"

	"go.uber.org/zap"
)

type Device interface {
	ID() DeviceID
	Close() error

	SetPropertyBool(property Property, value bool) error
	SetPropertyInt(property Property, value int) error
	SetPropertyDuration(property Property, value time.Duration) error

	InitiatorInit() error
	InitiatorSelectPassiveTarget(m Modulation, initData []byte) (*NfcTarget, error)
	InitiatorTransceiveBytes(tx, rx []byte, timeout time.Duration) (int, error)
	InitiatorDeselectTarget() error

	//WIP
	SetLastError(err error)
	GetInfiniteSelect() bool
	Logger() *zap.SugaredLogger
	DeviceGetSupportedModulation(mode Mode) ([]ModulationType, error)
	GetSupportedBaudRate(nmt ModulationType) ([]BaudRate, error)
	GetSupportedBaudRateTargetMode(nmt ModulationType) ([]BaudRate, error)
}

type NFCDeviceCommon struct {
	LastError      error
	InfiniteSelect bool
	Par            bool
	EasyFraming    bool
}
