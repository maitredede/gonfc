package pigpio

import (
	"context"
	"time"

	"github.com/maitredede/go-pigpiod"
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
	return ""
}

func (d *PN532PiGPIOI2CDeviceID) Open(logger *zap.SugaredLogger) (gonfc.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := pigpiod.Connect(ctx, d.drv.host)
	if err != nil {
		return nil, err
	}

	h, err := c.I2CO(d.drv.i2cBus, d.drv.i2cAddress, d.drv.i2cFlags)
	if err != nil {
		defer c.Close()
		return nil, err
	}

	dev := &PN532PiGPIOI2CDevice{
		id:     d,
		client: c,
		handle: h,
	}
	return dev, nil
}
