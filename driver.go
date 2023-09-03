package gonfc

import (
	"go.uber.org/zap"
)

type Driver interface {
	Manufacturer() string
	Product() string
	LookupDevices(logger *zap.SugaredLogger) ([]DeviceID, error)

	Conflicts(otherDriver Driver) bool
	String() string
}

type DeviceID interface {
	Driver() Driver
	Path() string
	Open(logger *zap.SugaredLogger) (Device, error)
	String() string
}
