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

	logger = log.Sugar()
	logger.Infof("gonfc version of nfc-list")

	drvs := common.RegisterAllDrivers(logger)

	devices := make([]gonfc.DeviceID, 0)
	for _, d := range drvs {
		dd, err := d.LookupDevices(logger)
		if err != nil {
			logger.Warnf("driver %v lookup error: %v", d, err)
			continue
		}
		logger.Infof("driver %v %v found %v devices", d.Manufacturer(), d.Product(), len(dd))
		devices = append(devices, dd...)
	}

	logger.Infof("found %v devices", len(devices))
	for _, dev := range devices {
		mainDev(dev)
	}
}

func mainDev(devID gonfc.DeviceID) {
	device, err := devID.Open(logger)
	if err != nil {
		logger.Warnf("device %v open error: %v", devID, err)
		return
	}
	defer device.Close()

	if err := device.InitiatorInit(); err != nil {
		logger.Errorf("initiator_init: %w", err)
		return
	}

	logger.Infof("NFC Device %v opened", device)

	listDev(device, gonfc.NMT_ISO14443A, gonfc.Nbr106)
}

func listDev(device gonfc.Device, modulationType gonfc.ModulationType, speed gonfc.BaudRate) {
	m := gonfc.Modulation{Type: modulationType, BaudRate: speed}
	targets, err := gonfc.InitiatorListPassiveTargets(device, m)
	if err != nil {
		logger.Warnf("device %v list error: %v", device.ID(), err)
		return
	}
	logger.Infof("  found %v targets", len(targets))
	for _, t := range targets {
		logger.Infof("  %v", t)
	}
}
