package pigpio

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/maitredede/go-pigpiod"
	"github.com/maitredede/gonfc"
)

type PN532PiGPIOI2CDriver struct {
	host       string
	i2cBus     byte
	i2cAddress byte
}

var (
	_ gonfc.Driver = (*PN532PiGPIOI2CDriver)(nil)
)

func NewDriver(host string, i2cBus byte, i2cAddress byte) (*PN532PiGPIOI2CDriver, error) {
	//TODO args validation
	drv := &PN532PiGPIOI2CDriver{
		host:       host,
		i2cBus:     i2cBus,
		i2cAddress: i2cAddress,
	}
	return drv, nil
}

func (d *PN532PiGPIOI2CDriver) Manufacturer() string {
	return "PN532 I2C via PiGPIO"
}

func (d *PN532PiGPIOI2CDriver) Product() string {
	return d.host
}

func (d *PN532PiGPIOI2CDriver) LookupDevices() ([]gonfc.DeviceID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := pigpiod.Connect(ctx, d.host)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	h, err := c.I2CO(d.i2cBus, d.i2cAddress, 0)
	if err != nil {
		return nil, err
	}
	defer c.I2CC(h)

	return nil, errors.New("TODO")
}

func (d *PN532PiGPIOI2CDriver) Conflicts(otherDriver gonfc.Driver) bool {
	if o, ok := otherDriver.(*PN532PiGPIOI2CDriver); ok {
		return strings.EqualFold(o.host, d.host)
	}
	return false
}
