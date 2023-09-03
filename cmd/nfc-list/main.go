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

	logger = log.Sugar()
	logger.Infof("gonfc version of nfc-list")

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

	logger.Debugf("=============================================")
	logger.Debugf("=== InitiatorInit")

	if err := gonfc.InitiatorInit(device); err != nil {
		logger.Errorf("initiator_init: %w", err)
		return
	}

	logger.Infof("NFC Device %v opened", device)

	listDev(device, gonfc.NMT_ISO14443A, gonfc.Nbr106)
	// listDev(device, gonfc.NMT_FELICA, gonfc.Nbr212)
	// listDev(device, gonfc.NMT_FELICA, gonfc.Nbr424)
	// listDev(device, gonfc.NMT_ISO14443B, gonfc.Nbr106)
	// listDev(device, gonfc.NMT_ISO14443BI, gonfc.Nbr106)
	// listDev(device, gonfc.NMT_ISO14443B2SR, gonfc.Nbr106)
	// listDev(device, gonfc.NMT_ISO14443B2CT, gonfc.Nbr106)
	// listDev(device, gonfc.NMT_ISO14443BICLASS, gonfc.Nbr106)
	// listDev(device, gonfc.NMT_JEWEL, gonfc.Nbr106)
	// listDev(device, gonfc.NMT_BARCODE, gonfc.Nbr106)
}

func listDev(device gonfc.Device, modulationType gonfc.ModulationType, speed gonfc.BaudRate) {
	logger.Debugf("=============================================")
	logger.Debugf("=== InitiatorListPassiveTargets")

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
