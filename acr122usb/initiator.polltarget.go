package acr122usb

import "github.com/maitredede/gonfc"

func (d *Acr122UsbDevice) InitiatorPollTarget(modulations []gonfc.Modulation, pollnr byte, period byte) (*gonfc.NfcTarget, error) {
	return d.chip.InitiatorPollTarget(modulations, pollnr, period)
}
