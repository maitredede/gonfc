package acr122usb

import (
	"time"

	"github.com/maitredede/gonfc"
)

func (pnd *Acr122UsbDevice) InitiatorInit() error {
	// Drop the field for a while
	if err := pnd.SetPropertyBool(gonfc.NP_ACTIVATE_FIELD, false); err != nil {
		return err
	}
	// Enable field so more power consuming cards can power themselves up
	if err := pnd.SetPropertyBool(gonfc.NP_ACTIVATE_FIELD, true); err != nil {
		return err
	}
	// Let the device try forever to find a target/tag
	if err := pnd.SetPropertyBool(gonfc.NP_INFINITE_SELECT, true); err != nil {
		return err
	}
	// Activate auto ISO14443-4 switching by default
	if err := pnd.SetPropertyBool(gonfc.NP_AUTO_ISO14443_4, true); err != nil {
		return err
	}
	// Force 14443-A mode
	if err := pnd.SetPropertyBool(gonfc.NP_FORCE_ISO14443_A, true); err != nil {
		return err
	}
	// Force speed at 106kbps
	if err := pnd.SetPropertyBool(gonfc.NP_FORCE_SPEED_106, true); err != nil {
		return err
	}
	// Disallow invalid frame
	if err := pnd.SetPropertyBool(gonfc.NP_ACCEPT_INVALID_FRAMES, false); err != nil {
		return err
	}
	// Disallow multiple frames
	if err := pnd.SetPropertyBool(gonfc.NP_ACCEPT_MULTIPLE_FRAMES, false); err != nil {
		return err
	}
	return pnd.chip.InitiatorInit()
}

func (pnd *Acr122UsbDevice) DeviceValidateModulation(mode gonfc.Mode, nm gonfc.Modulation) error {

	nmt, err := pnd.DeviceGetSupportedModulation(mode)
	if err != nil {
		return err
	}
	for _, i := range nmt {
		if i != nm.Type {
			continue
		}
		var nbr []gonfc.BaudRate
		var err error
		if mode == gonfc.N_INITIATOR {
			nbr, err = pnd.GetSupportedBaudRate(i)
			if err != nil {
				return err
			}
		} else {
			nbr, err = pnd.GetSupportedBaudRateTargetMode(i)
			if err != nil {
				return err
			}
		}
		for _, j := range nbr {
			if j == nm.BaudRate {
				return nil
			}
		}
		return gonfc.NFC_EINVARG
	}
	return gonfc.NFC_EINVARG
}

func (pnd *Acr122UsbDevice) InitiatorTransceiveBytes(tx, rx []byte, timeout time.Duration) (int, error) {
	return pnd.chip.InitiatorTransceiveBytes(tx, rx, timeout)
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
