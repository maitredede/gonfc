package pigpio

import (
	"context"
	"time"

	"github.com/maitredede/go-pigpiod"
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/compat"
	"github.com/maitredede/gonfc/pn53x"
	"go.uber.org/zap"
)

type PN532PiGPIOSPIDeviceID struct {
	drv *PN532PiGPIOSPIDriver
}

var _ (gonfc.DeviceID) = (*PN532PiGPIOSPIDeviceID)(nil)

func (d *PN532PiGPIOSPIDeviceID) Driver() gonfc.Driver {
	return d.drv
}

func (d *PN532PiGPIOSPIDeviceID) Path() string {
	return ""
}

func (d *PN532PiGPIOSPIDeviceID) Open(logger *zap.SugaredLogger) (gonfc.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := pigpiod.Connect(ctx, d.drv.host)
	if err != nil {
		return nil, err
	}

	dev := &PN532PiGPIOSPIDevice{
		id:     d,
		client: c,
	}
	io := &pn532spiIO{
		device: dev,
	}

	bInfiniteSelect := compat.NewBoolFieldGetSet(
		func() bool { return dev.InfiniteSelect },
		func(b bool) { dev.InfiniteSelect = b },
	)
	lastError := compat.NewErrorFieldGetSet(
		func() error { return dev.LastError },
		func(b error) { dev.LastError = b },
	)
	bPar := compat.NewBoolFieldGetSet(
		func() bool { return dev.Par },
		func(b bool) { dev.Par = b },
	)
	bEasyFraming := compat.NewBoolFieldGetSet(
		func() bool { return dev.EasyFraming },
		func(b bool) { dev.EasyFraming = b },
	)

	chip, err := pn53x.NewChip(io, logger.Named("pn53x"), bInfiniteSelect, lastError, bPar, bEasyFraming)
	if err != nil {
		defer c.Close()
		return nil, err
	}
	dev.chip = chip

	h, err := c.SPIO(d.drv.spiChannel, d.drv.spiBaud, d.drv.spiFlags)
	if err != nil {
		defer c.Close()
		return nil, err
	}
	dev.handle = h

	return dev, nil
}
