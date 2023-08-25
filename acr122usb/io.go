package acr122usb

import "github.com/maitredede/gonfc/pn53x"

type acr122io struct {
	dev *Acr122UsbDevice
}

var _ pn53x.IO = (*acr122io)(nil)

func (w *acr122io) Send(data []byte, timeout int) (int, error) {
	return w.dev.usbSend(data, timeout)
}

func (w *acr122io) Receive(data []byte, timeout int) (int, error) {
	return w.dev.usbReceive(data, timeout)
}

func (w *acr122io) MaxPacketSize() int {
	return w.dev.maxPacketSize
}
