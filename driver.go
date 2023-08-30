package gonfc

import (
	"go.uber.org/zap"
)

type Driver interface {
	Manufacturer() string
	Product() string
	LookupDevices(logger *zap.SugaredLogger) ([]DeviceID, error)

	Conflicts(otherDriver Driver) bool
}

type DeviceID interface {
	Driver() Driver
	Path() string
	Open(logger *zap.SugaredLogger) (Device, error)
}
