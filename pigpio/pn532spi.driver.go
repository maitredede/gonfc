package pigpio

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/maitredede/go-pigpiod"
	"github.com/maitredede/gonfc"
)

type PN532PiGPIOSPIDriver struct {
	host       string
	spiChannel byte
	spiBaud    uint32
	spiFlags   uint32
}

var (
	_ gonfc.Driver = (*PN532PiGPIOSPIDriver)(nil)
)

func NewSPIDriver(host string, spiChannel byte, spiBaud uint32, spiFlags uint32) (*PN532PiGPIOSPIDriver, error) {
	//TODO args validation
	drv := &PN532PiGPIOSPIDriver{
		host:       host,
		spiChannel: spiChannel,
		spiBaud:    spiBaud,
		spiFlags:   spiFlags,
	}
	return drv, nil
}

func (d *PN532PiGPIOSPIDriver) Manufacturer() string {
	return "PN532 SPI via PiGPIO"
}

func (d *PN532PiGPIOSPIDriver) Product() string {
	return d.host
}

func (d *PN532PiGPIOSPIDriver) LookupDevices() ([]gonfc.DeviceID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := pigpiod.Connect(ctx, d.host)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	h, err := c.SPIO(d.spiChannel, d.spiBaud, d.spiFlags)
	if err != nil {
		return nil, err
	}
	defer c.SPIC(h)

	return nil, errors.New("TODO")
}

func (d *PN532PiGPIOSPIDriver) Conflicts(otherDriver gonfc.Driver) bool {
	if o, ok := otherDriver.(*PN532PiGPIOSPIDriver); ok {
		return strings.EqualFold(o.host, d.host)
	}
	return false
}
