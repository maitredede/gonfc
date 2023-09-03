package periphio

import (
	"fmt"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/compat"
	"github.com/maitredede/gonfc/pn53x"
	"go.uber.org/zap"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

type PeriphioI2CDriver struct {
	busName string
	addr    uint16
}

var _ gonfc.Driver = (*PeriphioI2CDriver)(nil)

func NewPeriphioI2CDriver(busName string, address uint16) (*PeriphioI2CDriver, error) {
	// Make sure periph is initialized.
	// TODO: Use host.Init(). It is not used in this example to prevent circular
	// go package import.
	if _, err := driverreg.Init(); err != nil {
		return nil, fmt.Errorf("periph.io driverreg.Init() error: %w", err)
	}

	d := PeriphioI2CDriver{
		busName: busName,
		addr:    address,
	}
	return &d, nil
}

func (d *PeriphioI2CDriver) Manufacturer() string {
	return "gonfc/periphio"
}

func (d *PeriphioI2CDriver) Product() string {
	return "pn532-i2c"
}

func (d *PeriphioI2CDriver) String() string {
	return fmt.Sprintf("%s %s", d.Manufacturer(), d.Product())
}

func (d *PeriphioI2CDriver) LookupDevices(logger *zap.SugaredLogger) ([]gonfc.DeviceID, error) {
	dev, err := d.openDevice(logger)
	if err != nil {
		return nil, err
	}
	defer dev.Close()

	id := &PeriphioI2CDeviceID{
		drv: d,
	}
	return []gonfc.DeviceID{id}, nil
}

func (d *PeriphioI2CDriver) openDevice(logger *zap.SugaredLogger) (*PeriphioI2CDevice, error) {
	// Use i2creg I²C hwBus registry to find the first available I²C hwBus.
	hwBus, err := i2creg.Open(d.busName)
	if err != nil {
		return nil, fmt.Errorf("i2c bus open failed: %w", err)
	}
	// Dev is a valid conn.Conn.
	hwDev := &i2c.Dev{Addr: d.addr, Bus: hwBus}

	dev := &PeriphioI2CDevice{
		driver: d,
		logger: logger,
		hwBus:  hwBus,
		hwDev:  hwDev,
	}
	io := &periphioI2Cio{
		dev: dev,
	}
	dev.io = io

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

	return dev, nil
}

func (d *PeriphioI2CDriver) Conflicts(otherDriver gonfc.Driver) bool {
	other, ok := otherDriver.(*PeriphioI2CDriver)
	if !ok {
		return false
	}
	return d.busName == other.busName && d.addr == other.addr
}
