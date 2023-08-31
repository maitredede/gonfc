package acr122usb

import (
	"time"

	"github.com/maitredede/gonfc"
)

func (pnd *Acr122UsbDevice) InitiatorInit() error {
	return pnd.chip.InitiatorInit()
}

func (pnd *Acr122UsbDevice) InitiatorTransceiveBytes(tx, rx []byte, timeout time.Duration) (int, error) {
	return pnd.chip.InitiatorTransceiveBytes(tx, rx, timeout)
}

func (d *Acr122UsbDevice) InitiatorSelectPassiveTarget(m gonfc.Modulation, initData []byte) (*gonfc.NfcTarget, error) {
	return d.chip.InitiatorSelectPassiveTarget(m, initData)
}

func (pnd *Acr122UsbDevice) InitiatorDeselectTarget() error {
	return pnd.chip.InitiatorDeselectTarget()
}

func (pnd *Acr122UsbDevice) DeviceGetSupportedModulation(mode gonfc.Mode) ([]gonfc.ModulationType, error) {
	return pnd.chip.GetSupportedModulation(mode)
}

func (pnd *Acr122UsbDevice) GetSupportedBaudRate(mt gonfc.ModulationType) ([]gonfc.BaudRate, error) {
	return pnd.chip.GetSupportedBaudRate(gonfc.N_INITIATOR, mt)
}

func (pnd *Acr122UsbDevice) GetSupportedBaudRateTargetMode(mt gonfc.ModulationType) ([]gonfc.BaudRate, error) {
	return pnd.chip.GetSupportedBaudRate(gonfc.N_TARGET, mt)
}
