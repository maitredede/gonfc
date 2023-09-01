package main

import (
	goflag "flag"
	"time"

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
	logger.Infof("gonfc version of nfc-poll")

	//periphio
	_, err := host.Init()
	if err != nil {
		logger.Fatal(err)
	}

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

	nmModulations := []gonfc.Modulation{
		{Type: gonfc.NMT_ISO14443A, BaudRate: gonfc.Nbr106},
		{Type: gonfc.NMT_ISO14443B, BaudRate: gonfc.Nbr106},
		{Type: gonfc.NMT_FELICA, BaudRate: gonfc.Nbr212},
		{Type: gonfc.NMT_FELICA, BaudRate: gonfc.Nbr424},
		{Type: gonfc.NMT_JEWEL, BaudRate: gonfc.Nbr106},
		{Type: gonfc.NMT_ISO14443BICLASS, BaudRate: gonfc.Nbr106},
	}
	szModulations := len(nmModulations)
	uiPollNr := 20
	uiPeriod := 2
	uiPeriodDuration := time.Duration(uiPeriod) * 150 * time.Millisecond
	totalDuration := uiPeriodDuration * time.Duration(uiPollNr) * time.Duration(szModulations)

	logger.Infof("NFC reader: %s opened", device)
	logger.Infof("NFC device will poll during %v (%v pollings of %v for %v modulations)", totalDuration, uiPollNr, uiPeriodDuration, szModulations)

	if err := device.InitiatorInit(); err != nil {
		logger.Errorf("initiator_init failed: %v", err)
	}
	nt, err := gonfc.InitiatorPollTarget(device, nmModulations, uiPollNr, uiPeriod)
	if err != nil {
		logger.Errorf("poll target error: %v", err)
		return
	}

	if nt == nil {
		logger.Info("No target found.")
		return
	}

	logger.Infof("Waiting for card removing...")
	for {
		isPresent, err := gonfc.InitiatorTargetIsPresent(device /*nt ? */, nil)
		if err != nil {
			logger.Errorf("target is present error: %v", err)
			return
		}
		if !isPresent {
			logger.Infof("done.")
			break
		}
	}
}
