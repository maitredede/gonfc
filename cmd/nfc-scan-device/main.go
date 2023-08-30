package main

import (
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/acr122usb"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	redir := zap.RedirectStdLog(log)
	defer redir()

	logger := log.Sugar()

	drvs := []gonfc.Driver{
		&acr122usb.Acr122USBDriver{},
	}

	devices := make([]gonfc.DeviceID, 0)
	for _, d := range drvs {
		dd, err := d.LookupDevices(logger)
		if err != nil {
			continue
		}
		devices = append(devices, dd...)
	}

	logger.Infof("found %v devices", len(devices))
	for _, dev := range devices {
		logger.Infof("  %v %v: %v", dev.Driver().Manufacturer(), dev.Driver().Product(), dev.Path())
	}
}
