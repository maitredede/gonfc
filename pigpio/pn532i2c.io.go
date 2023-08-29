package pigpio

import "github.com/maitredede/gonfc/pn53x"

type pn532i2cIO struct {
	device *PN532PiGPIOI2CDevice
}

var _ pn53x.IO = (*pn532i2cIO)(nil)

func (d *pn532i2cIO) Send(data []byte, timeout int) (int, error) {
	// d.device.client.I2CW()
	panic("TODO")
}

func (d *pn532i2cIO) Receive(data []byte, timeout int) (int, error) {
	// d.device.client.I2CR()
	panic("TODO")
}

func (d *pn532i2cIO) MaxPacketSize() int {
	panic("TODO")
}
