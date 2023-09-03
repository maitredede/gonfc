package acr122usb

import (
	"fmt"
	"slices"

	"github.com/google/gousb"
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/compat"
	"github.com/maitredede/gonfc/pn53x"
	"go.uber.org/zap"
)

type acr122DeviceID struct {
	driver     *Acr122USBDriver
	desc       *gousb.DeviceDesc
	deviceInfo Acr122UsbSupportedDevice
	uif        gousb.InterfaceSetting
}

var _ gonfc.DeviceID = (*acr122DeviceID)(nil)

func (d *acr122DeviceID) Driver() gonfc.Driver {
	return d.driver
}

func (d *acr122DeviceID) Path() string {
	return d.desc.String()
}

func (d *acr122DeviceID) String() string {
	return d.deviceInfo.Name + " " + d.desc.String()
}

// Open opens an acr122 usb device (libnfc: acr122_usb_open)
func (d *acr122DeviceID) Open(logger *zap.SugaredLogger) (gonfc.Device, error) {
	devs, err := d.driver.usb.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return slices.Equal(desc.Path, d.desc.Path)
	})
	if len(devs) == 0 {
		return nil, fmt.Errorf("device not found")
	}
	if len(devs) > 1 || err != nil {
		defer func() {
			for _, d := range devs {
				d.Close()
			}
		}()
		if err == nil {
			return nil, fmt.Errorf("multiple devices found")
		}
		return nil, err
	}

	dev := &Acr122UsbDevice{
		id:     d,
		device: devs[0],
		logger: logger,
	}
	dev.io = &acr122io{
		dev: dev,
	}

	if err := dev.device.Reset(); err != nil {
		defer dev.Close()
		return nil, err
	}

	// libnfc.usb_claim_interface()
	intf, done, err := dev.device.DefaultInterface()
	if err != nil {
		defer dev.Close()
		return nil, err
	}
	dev.intf = intf
	dev.done = done

	// libnfc.acr122_usb_get_end_points()
	// 3 Endpoints maximum: Interrupt In, Bulk In, Bulk Out
	if err := dev.usbGetEndPoints(); err != nil {
		defer dev.Close()
		return nil, err
	}

	// libnfc.acr122_usb_get_usb_device_name
	dev.name = dev.usbGetUsbDeviceName()
	abortFlag := compat.NewBoolFieldGetSet(func() bool { return dev.AbortFlag }, func(b bool) { dev.AbortFlag = b })
	bInfiniteSelect := compat.NewBoolFieldGetSet(func() bool { return dev.InfiniteSelect }, func(b bool) { dev.InfiniteSelect = b })
	lastError := compat.NewErrorFieldGetSet(func() error { return dev.LastError }, func(b error) { dev.LastError = b })
	bPar := compat.NewBoolFieldGetSet(func() bool { return dev.Par }, func(b bool) { dev.Par = b })
	bEasyFraming := compat.NewBoolFieldGetSet(func() bool { return dev.EasyFraming }, func(b bool) { dev.EasyFraming = b })
	chip, err := pn53x.NewChip(
		dev.io,
		dev.logger.Named("pn53x"),
		bInfiniteSelect,
		lastError,
		bPar,
		bEasyFraming,
		abortFlag,
	)
	if err != nil {
		dev.Close()
		return nil, err
	}
	dev.chip = chip

	dev.timerCorrection = 46 // empirical tuning

	if err := dev.usbinit(); err != nil {
		defer dev.Close()
		return nil, fmt.Errorf("usbinit: %w", err)
	}

	logger.Infof("blah blah blah has been claimed")

	return dev, nil
}
