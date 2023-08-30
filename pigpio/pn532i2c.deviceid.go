package pigpio

import (
	"fmt"

	"github.com/maitredede/gonfc"
	"go.uber.org/zap"
)

type PN532PiGPIOI2CDeviceID struct {
	drv *PN532PiGPIOI2CDriver
}

var _ (gonfc.DeviceID) = (*PN532PiGPIOI2CDeviceID)(nil)

func (d *PN532PiGPIOI2CDeviceID) Driver() gonfc.Driver {
	return d.drv
}

func (d *PN532PiGPIOI2CDeviceID) Path() string {
	return fmt.Sprintf("tcp://%v?bus=%v&addr=%v&flags=%v", d.drv.host, d.drv.i2cBus, d.drv.i2cAddress, d.drv.i2cFlags)
}

func (d *PN532PiGPIOI2CDeviceID) Open(logger *zap.SugaredLogger) (gonfc.Device, error) {
	dev, err := d.drv.openDevice(logger)
	if err != nil {
		return nil, err
	}
	if err := dev.chip.Init(); err != nil {
		defer dev.Close()
		return nil, err
	}
	return dev, nil
}
