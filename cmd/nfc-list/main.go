package main

import (
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/acr122usb"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	redir := zap.RedirectStdLog(log)
	defer redir()

	logger = log.Sugar()

	drvs := []gonfc.Driver{
		&acr122usb.Acr122USBDriver{},
	}

	devices := make([]gonfc.DeviceID, 0)
	for _, d := range drvs {
		dd, err := d.LookupDevices()
		if err != nil {
			continue
		}
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
	}
	defer device.Close()

	logger.Debugf("=== device opened ===")

	if err := device.InitiatorInit(); err != nil {
		logger.Errorf("initiator_init: %w", err)
		return
	}

	logger.Infof("NFC Device %v opened", device)
	logger.Infof("=== nfc_initiator_init ok ===")

	listDev(device, gonfc.NMT_ISO14443A, gonfc.Nbr106)
}

func listDev(device gonfc.Device, modulationType gonfc.ModulationType, speed gonfc.BaudRate) {
	m := gonfc.Modulation{Type: modulationType, BaudRate: speed}
	targets, err := device.InitiatorListPassiveTargets(m)
	if err != nil {
		logger.Warnf("device %v list error: %v", device.ID(), err)
		return
	}
	logger.Infof("  found %v targets", len(targets))
	for _, t := range targets {
		logger.Infof("  %v", t)
	}
}
