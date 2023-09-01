package periphio

import (
	"fmt"

	"github.com/maitredede/gonfc"
	"go.uber.org/zap"
)

type PeriphioI2CDeviceID struct {
	drv *PeriphioI2CDriver
}

var _ gonfc.DeviceID = (*PeriphioI2CDeviceID)(nil)

func (d *PeriphioI2CDeviceID) Driver() gonfc.Driver {
	return d.drv
}

func (d *PeriphioI2CDeviceID) Path() string {
	return fmt.Sprintf("periph.io://i2c/%v/%v", d.drv.busName, d.drv.addr)
}

func (d *PeriphioI2CDeviceID) Open(logger *zap.SugaredLogger) (gonfc.Device, error) {
	return d.drv.openDevice(logger)
}
