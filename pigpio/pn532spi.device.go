package pigpio

import (
	"errors"

	"github.com/maitredede/go-pigpiod"
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/pn53x"
)

type PN532PiGPIOSPIDevice struct {
	id     *PN532PiGPIOSPIDeviceID
	client *pigpiod.Conn
	handle uint32

	chip *pn53x.Chip

	gonfc.NFCDeviceCommon
}

var _ gonfc.Device = (*PN532PiGPIOSPIDevice)(nil)

func (d *PN532PiGPIOSPIDevice) ID() gonfc.DeviceID {
	return d.id
}

func (d *PN532PiGPIOSPIDevice) Close() error {
	errs := make([]error, 0)
	_, err := d.client.SPIC(d.handle)
	errs = append(errs, err)
	err = d.client.Close()
	errs = append(errs, err)
	return errors.Join(errs...)
}

func (d *PN532PiGPIOSPIDevice) SetPropertyBool(property gonfc.Property, value bool) error {
	panic("TODO")
}

func (d *PN532PiGPIOSPIDevice) SetPropertyInt(property gonfc.Property, value int) error {
	panic("TODO")
}

func (d *PN532PiGPIOSPIDevice) InitiatorInit() error {
	panic("TODO")
}

func (d *PN532PiGPIOSPIDevice) InitiatorListPassiveTargets(m gonfc.Modulation) ([]*gonfc.NfcTarget, error) {
	panic("TODO")
}

func (d *PN532PiGPIOSPIDevice) InitiatorSelectPassiveTarget(m gonfc.Modulation, initData []byte) (*gonfc.NfcTarget, error) {
	panic("TODO")
}

func (d *PN532PiGPIOSPIDevice) InitiatorTransceiveBytes(tx, rx []byte, timeout int) (int, error) {
	panic("TODO")
}

func (d *PN532PiGPIOSPIDevice) InitiatorDeselectTarget() error {
	panic("TODO")
}
