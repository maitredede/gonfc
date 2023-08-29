package gonfc

import "go.uber.org/zap"

type Driver interface {
	Manufacturer() string
	Product() string
	LookupDevices() ([]DeviceID, error)

	Conflicts(otherDriver Driver) bool
}

type DeviceID interface {
	Driver() Driver
	Path() string
	Open(logger *zap.SugaredLogger) (Device, error)
}

type Device interface {
	ID() DeviceID
	Close() error

	SetPropertyBool(property Property, value bool) error
	SetPropertyInt(property Property, value int) error

	InitiatorInit() error
	InitiatorListPassiveTargets(m Modulation) ([]*NfcTarget, error)
	InitiatorSelectPassiveTarget(m Modulation, initData []byte) (*NfcTarget, error)
	InitiatorTransceiveBytes(tx, rx []byte, timeout int) (int, error)
	InitiatorDeselectTarget() error
}

type NFCDeviceCommon struct {
	LastError      error
	InfiniteSelect bool
	Par            bool
	EasyFraming    bool
}
