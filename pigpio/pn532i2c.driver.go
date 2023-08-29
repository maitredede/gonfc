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
	i2cFlags   uint32
}

var (
	_ gonfc.Driver = (*PN532PiGPIOI2CDriver)(nil)
)

func NewI2CDriver(host string, i2cBus byte, i2cAddress byte, i2cFlags uint32) (*PN532PiGPIOI2CDriver, error) {
	//TODO args validation
	drv := &PN532PiGPIOI2CDriver{
		host:       host,
		i2cBus:     i2cBus,
		i2cAddress: i2cAddress,
		i2cFlags:   i2cFlags,
	}
	return drv, nil
}

func (d *PN532PiGPIOI2CDriver) Manufacturer() string {
	return "PN532 I2C via PiGPIO"
}

func (d *PN532PiGPIOI2CDriver) Product() string {
	return d.host
}

func (d *PN532PiGPIOI2CDriver) Conflicts(otherDriver gonfc.Driver) bool {
	if o, ok := otherDriver.(*PN532PiGPIOI2CDriver); ok {
		return strings.EqualFold(o.host, d.host)
	}
	return false
}

func (d *PN532PiGPIOI2CDriver) LookupDevices() ([]gonfc.DeviceID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := pigpiod.Connect(ctx, d.host)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	h, err := c.I2CO(d.i2cBus, d.i2cAddress, d.i2cFlags)
	if err != nil {
		return nil, err
	}
	defer c.I2CC(h)

	return nil, errors.New("TODO")
}
