package pigpio

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/maitredede/go-pigpiod"
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/compat"
	"github.com/maitredede/gonfc/pn53x"
	"go.uber.org/zap"
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
	return "PiGPIO - PN532 I2C"
}

func (d *PN532PiGPIOI2CDriver) Product() string {
	return d.host
}

func (d *PN532PiGPIOI2CDriver) String() string {
	return fmt.Sprintf("%s %s", d.Manufacturer(), d.Product())
}

func (d *PN532PiGPIOI2CDriver) Conflicts(otherDriver gonfc.Driver) bool {
	if o, ok := otherDriver.(*PN532PiGPIOI2CDriver); ok {
		return strings.EqualFold(o.host, d.host)
	}
	return false
}

func (d *PN532PiGPIOI2CDriver) LookupDevices(logger *zap.SugaredLogger) ([]gonfc.DeviceID, error) {
	dev, err := d.openDevice(logger)
	if dev != nil {
		defer dev.Close()
	}
	if err != nil {
		return nil, err
	}
	return []gonfc.DeviceID{dev.id}, nil
}

func (d *PN532PiGPIOI2CDriver) openDevice(logger *zap.SugaredLogger) (*PN532PiGPIOI2CDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := pigpiod.Connect(ctx, d.host)
	if err != nil {
		return nil, fmt.Errorf("tcp connect error: %w", err)
	}

	h, err := c.I2CO(d.i2cBus, d.i2cAddress, d.i2cFlags)
	if err != nil {
		defer c.Close()
		return nil, fmt.Errorf("i2c open error: %w", err)
	}

	dev := &PN532PiGPIOI2CDevice{
		id: &PN532PiGPIOI2CDeviceID{
			drv: d,
		},
		NFCDeviceCommon: gonfc.NFCDeviceCommon{},
		client:          c,
		handle:          h,
		logger:          logger,
	}
	io := &pn532i2cIO{
		device: dev,
	}
	abortFlag := compat.NewBoolFieldGetSet(func() bool { return dev.AbortFlag }, func(b bool) { dev.AbortFlag = b })
	lastError := compat.NewErrorFieldGetSet(func() error { return dev.LastError }, func(err error) { dev.LastError = err })
	bInfiniteSelect := compat.NewBoolFieldGetSet(func() bool { return dev.InfiniteSelect }, func(b bool) { dev.InfiniteSelect = b })
	bPar := compat.NewBoolFieldGetSet(func() bool { return dev.Par }, func(b bool) { dev.Par = b })
	bEasyFraming := compat.NewBoolFieldGetSet(func() bool { return dev.EasyFraming }, func(b bool) { dev.EasyFraming = b })
	chip, err := pn53x.NewChipPN532(
		io,
		logger,
		bInfiniteSelect,
		lastError,
		bPar,
		bEasyFraming,
		abortFlag,
	)
	if err != nil {
		defer dev.Close()
		return nil, err
	}
	dev.chip = chip

	if err := chip.CheckCommunication(); err != nil {
		defer dev.Close()
		return nil, err
	}

	return dev, nil
}
