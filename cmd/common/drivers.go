package common

import (
	"flag"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/acr122usb"
	"github.com/maitredede/gonfc/periphio"
	"github.com/maitredede/gonfc/pigpio"
	"go.uber.org/zap"
)

var (
	disablePigpioI2C   bool
	disablePeriphioI2C bool
)

func init() {
	flag.BoolVar(&disablePigpioI2C, "disable-pidpio-i2c", false, "Disable pigpio I2C")
	flag.BoolVar(&disablePeriphioI2C, "disable-periphio-i2c", false, "Disable periph.io I2C")
}

func RegisterAllDrivers(logger *zap.SugaredLogger) []gonfc.Driver {
	drvs := []gonfc.Driver{
		&acr122usb.Acr122USBDriver{},
	}
	if !disablePigpioI2C {
		// Remote raspberry pi with pigpiod installed
		i2cDrv, err := pigpio.NewI2CDriver("10.105.3.76:8888", 1, 0x24, 0)
		if err != nil {
			logger.Warnf("pigpio i2c driver not available: %v", err)
		} else {
			drvs = append(drvs, i2cDrv)
		}
	}
	if !disablePeriphioI2C {
		// Local I2C using Periph.io
		i2cDrv, err := periphio.NewPeriphioI2CDriver("", 0x24)
		if err != nil {
			logger.Warnf("periphio i2c driver not available: %v", err)
		} else {
			drvs = append(drvs, i2cDrv)
		}
	}
	return drvs
}
