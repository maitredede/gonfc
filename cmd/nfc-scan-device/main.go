package main

import (
	"flag"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/cmd/common"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func main() {
	flag.Parse()

	log := common.InitLogger(true)
	defer log.Sync()
	redir := zap.RedirectStdLog(log)
	defer redir()

	logger := log.Sugar()
	logger.Infof("gonfc version of nfc-list")

	drvs := common.RegisterAllDrivers(logger)

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
