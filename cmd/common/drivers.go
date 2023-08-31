package common

import (
	"flag"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/acr122usb"
	"github.com/maitredede/gonfc/pigpio"
	"go.uber.org/zap"
)

var (
	disableRemoteI2C bool
)

func init() {
	flag.BoolVar(&disableRemoteI2C, "disable-remote-i2c", false, "Disable remote I2C")
}

func RegisterAllDrivers(logger *zap.SugaredLogger) []gonfc.Driver {
	drvs := []gonfc.Driver{
		&acr122usb.Acr122USBDriver{},
	}
	if !disableRemoteI2C {
		// I have a raspberry pi with pigpiod installed
		i2cDrv, err := pigpio.NewI2CDriver("10.105.3.76:8888", 1, 0x24, 0)
		if err != nil {
			logger.Warnf("i2c driver not available: %v", err)
		} else {
			drvs = append(drvs, i2cDrv)
		}
	}
	return drvs
}
