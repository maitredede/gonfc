package acr122usb

import (
	"sync"

	"github.com/google/gousb"
	"github.com/maitredede/gonfc"
)

var (
	usbctx = sync.OnceValue(initCtx)
)

func initCtx() *gousb.Context {
	return gousb.NewContext()
}

type Acr122USBDriver struct {
}

var _ gonfc.Driver = (*Acr122USBDriver)(nil)

func (Acr122USBDriver) Manufacturer() string {
	return "libusb"
}
func (Acr122USBDriver) Product() string {
	return "acr122"
}

func (Acr122USBDriver) Conflicts(otherDriver gonfc.Driver) bool {
	if _, ok := otherDriver.(*Acr122USBDriver); ok {
		return true
	}
	return false
}

func (d *Acr122USBDriver) LookupDevices() ([]gonfc.DeviceID, error) {

	c := usbctx()

	result := make([]gonfc.DeviceID, 0)

	devs, err := c.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		ok := false
		var deviceInfo Acr122UsbSupportedDevice
		for _, d := range UsbSupportedDevices {
			if d.VID == desc.Vendor && d.PID == desc.Product {
				ok = true
				deviceInfo = d
				break
			}
		}
		if !ok {
			return false
		}

		// checks from libnfc
		// Make sure there are 2 endpoints available
		// with libusb-win32 we got some null pointers so be robust before looking at endpoints:
		if len(desc.Configs) == 0 {
			return false
		}
		cfg := desc.Configs[1]
		if len(cfg.Interfaces) == 0 {
			return false
		}
		iface := cfg.Interfaces[0]
		if len(iface.AltSettings) == 0 {
			return false
		}
		aset := iface.AltSettings[0]
		if len(aset.Endpoints) < 2 {
			return false
		}

		id := &acr122DeviceID{
			driver:     d,
			desc:       desc,
			deviceInfo: deviceInfo,
			uif:        aset,
		}
		result = append(result, id)
		return false
	})

	for _, d := range devs {
		d.Close()
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
