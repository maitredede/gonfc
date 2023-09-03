package main

import (
	goflag "flag"

	"github.com/google/gousb"
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/cmd/common"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"periph.io/x/host/v3"
)

var logger *zap.SugaredLogger

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	log := common.InitLogger(true)
	defer log.Sync()
	redir := zap.RedirectStdLog(log)
	defer redir()

	logger := log.Sugar()
	logger.Infof("gonfc version of nfc-scan-device")

	//gousb
	usb := gousb.NewContext()
	defer usb.Close()

	//periphio
	_, err := host.Init()
	if err != nil {
		logger.Fatal(err)
	}

	drvs := common.RegisterAllDrivers(logger, usb)

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
