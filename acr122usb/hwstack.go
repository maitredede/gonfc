package acr122usb

// import (
// 	"fmt"
// 	"slices"
// 	"sync"

// 	"github.com/google/gousb"
// 	"go.uber.org/zap"
// 	"isi.nc/ddaly/hardware/hardware"
// 	"isi.nc/ddaly/hardware/nfc"
// 	"isi.nc/ddaly/hardware/nfc/gonfc/pn53x"
// )

// var (
// 	usbctx = sync.OnceValue(initCtx)
// )

// func initCtx() *gousb.Context {
// 	return gousb.NewContext()
// }

// type Acr122USBDriver struct {
// }

// var _ nfc.NFCDriver = (*Acr122USBDriver)(nil)

// func (Acr122USBDriver) Manufacturer() string {
// 	return "libusb"
// }
// func (Acr122USBDriver) Product() string {
// 	return "acr122"
// }

// func (d *Acr122USBDriver) LookupDevices() ([]hardware.DeviceID, error) {

// 	c := usbctx()

// 	result := make([]hardware.DeviceID, 0)

// 	devs, err := c.OpenDevices(func(desc *gousb.DeviceDesc) bool {
// 		ok := false
// 		var deviceInfo Acr122UsbSupportedDevice
// 		for _, d := range UsbSupportedDevices {
// 			if d.VID == desc.Vendor && d.PID == desc.Product {
// 				ok = true
// 				deviceInfo = d
// 				break
// 			}
// 		}
// 		if !ok {
// 			return false
// 		}

// 		// checks from libnfc
// 		// Make sure there are 2 endpoints available
// 		// with libusb-win32 we got some null pointers so be robust before looking at endpoints:
// 		if len(desc.Configs) == 0 {
// 			return false
// 		}
// 		cfg := desc.Configs[1]
// 		if len(cfg.Interfaces) == 0 {
// 			return false
// 		}
// 		iface := cfg.Interfaces[0]
// 		if len(iface.AltSettings) == 0 {
// 			return false
// 		}
// 		aset := iface.AltSettings[0]
// 		if len(aset.Endpoints) < 2 {
// 			return false
// 		}

// 		id := &acr122DeviceID{
// 			driver:     d,
// 			desc:       desc,
// 			deviceInfo: deviceInfo,
// 			uif:        aset,
// 		}
// 		result = append(result, id)
// 		return false
// 	})

// 	for _, d := range devs {
// 		d.Close()
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// type acr122DeviceID struct {
// 	driver     *Acr122USBDriver
// 	desc       *gousb.DeviceDesc
// 	deviceInfo Acr122UsbSupportedDevice
// 	uif        gousb.InterfaceSetting
// }

// var _ nfc.NFCDeviceID = (*acr122DeviceID)(nil)

// func (d *acr122DeviceID) Driver() hardware.Driver {
// 	return d.driver
// }

// func (d *acr122DeviceID) Path() string {
// 	return d.desc.String()
// }

// func (d *acr122DeviceID) String() string {
// 	return d.deviceInfo.Name + " " + d.desc.String()
// }

// // Open opens an acr122 usb device (libnfc: acr122_usb_open)
// func (d *acr122DeviceID) Open(logger *zap.SugaredLogger) (hardware.Device, error) {
// 	c := usbctx()

// 	devs, err := c.OpenDevices(func(desc *gousb.DeviceDesc) bool {
// 		return slices.Equal(desc.Path, d.desc.Path)
// 	})
// 	if len(devs) == 0 {
// 		return nil, fmt.Errorf("device not found")
// 	}
// 	if len(devs) > 1 || err != nil {
// 		defer func() {
// 			for _, d := range devs {
// 				d.Close()
// 			}
// 		}()
// 		if err == nil {
// 			return nil, fmt.Errorf("multiple devices found")
// 		}
// 		return nil, err
// 	}

// 	dev := &Acr122UsbDevice{
// 		id:     d,
// 		device: devs[0],
// 		logger: logger,
// 	}
// 	dev.io = &acr122io{
// 		dev: dev,
// 	}

// 	if err := dev.device.Reset(); err != nil {
// 		defer dev.Close()
// 		return nil, err
// 	}

// 	// libnfc.usb_claim_interface()
// 	intf, done, err := dev.device.DefaultInterface()
// 	if err != nil {
// 		defer dev.Close()
// 		return nil, err
// 	}
// 	dev.intf = intf
// 	dev.done = done

// 	// libnfc.acr122_usb_get_end_points()
// 	// 3 Endpoints maximum: Interrupt In, Bulk In, Bulk Out
// 	for adr, ep := range dev.id.uif.Endpoints {
// 		if ep.TransferType != gousb.TransferTypeBulk {
// 			continue
// 		}
// 		if ep.Direction == gousb.EndpointDirectionIn {
// 			dev.epInAdr = adr
// 			iep, err := dev.intf.InEndpoint(int(ep.Address))
// 			if err != nil {
// 				defer dev.Close()
// 				return nil, err
// 			}
// 			dev.epIn = iep
// 			dev.maxPacketSize = ep.MaxPacketSize
// 		}
// 		if ep.Direction == gousb.EndpointDirectionOut {
// 			dev.epOutAdr = adr
// 			oep, err := dev.intf.OutEndpoint(int(ep.Address))
// 			if err != nil {
// 				defer dev.Close()
// 				return nil, err
// 			}
// 			dev.epOut = oep
// 			dev.maxPacketSize = ep.MaxPacketSize
// 		}
// 	}

// 	// libnfc.acr122_usb_get_usb_device_name
// 	dev.name = dev.usbGetUsbDeviceName()

// 	chip, err := pn53x.NewChip(dev.io, dev.logger.Named("pn53x"))
// 	if err != nil {
// 		dev.Close()
// 		return nil, err
// 	}
// 	dev.chip = chip

// 	dev.timerCorrection = 46 // empirical tuning

// 	if err := dev.usbinit(); err != nil {
// 		defer dev.Close()
// 		return nil, fmt.Errorf("usbinit: %w", err)
// 	}

// 	return dev, nil
// }
