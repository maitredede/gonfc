package pigpio

import (
	"errors"

	"github.com/maitredede/go-pigpiod"
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/pn53x"
)

type PN532PiGPIOI2CDevice struct {
	id     *PN532PiGPIOI2CDeviceID
	client *pigpiod.Conn
	handle uint32

	chip *pn53x.Chip

	gonfc.NFCDeviceCommon
}

var _ gonfc.Device = (*PN532PiGPIOI2CDevice)(nil)

func (d *PN532PiGPIOI2CDevice) ID() gonfc.DeviceID {
	return d.id
}

func (d *PN532PiGPIOI2CDevice) Close() error {
	errs := make([]error, 0)
	_, err := d.client.I2CC(d.handle)
	errs = append(errs, err)
	err = d.client.Close()
	errs = append(errs, err)
	return errors.Join(errs...)
}

func (d *PN532PiGPIOI2CDevice) SetPropertyBool(property gonfc.Property, value bool) error {
	panic("TODO")
}

func (d *PN532PiGPIOI2CDevice) SetPropertyInt(property gonfc.Property, value int) error {
	panic("TODO")
}

func (d *PN532PiGPIOI2CDevice) InitiatorInit() error {
	panic("TODO")
}

func (d *PN532PiGPIOI2CDevice) InitiatorListPassiveTargets(m gonfc.Modulation) ([]*gonfc.NfcTarget, error) {
	panic("TODO")
}

func (d *PN532PiGPIOI2CDevice) InitiatorSelectPassiveTarget(m gonfc.Modulation, initData []byte) (*gonfc.NfcTarget, error) {
	panic("TODO")
}

func (d *PN532PiGPIOI2CDevice) InitiatorTransceiveBytes(tx, rx []byte, timeout int) (int, error) {
	panic("TODO")
}

func (d *PN532PiGPIOI2CDevice) InitiatorDeselectTarget() error {
	panic("TODO")
}
