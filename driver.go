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
	InitiatorListPassiveTargets(m Modulation) ([]Target, error)
	InitiatorSelectPassiveTarget(m Modulation, initData []byte) (Target, error)
	InitiatorTransceiveBytes(tx, rx []byte, timeout int) (int, error)
	InitiatorDeselectTarget() error
}

type NFCDeviceCommon struct {
	InfiniteSelect bool
	Par            bool
	EasyFraming    bool
}