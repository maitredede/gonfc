package pigpio

import "github.com/maitredede/gonfc/pn53x"

type pn532spiIO struct {
	device *PN532PiGPIOSPIDevice
}

var _ pn53x.IO = (*pn532spiIO)(nil)

func (d *pn532spiIO) Send(data []byte, timeout int) (int, error) {
	err := d.device.client.SPIW(d.device.handle, data)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (d *pn532spiIO) Receive(data []byte, timeout int) (int, error) {
	num := uint32(len(data))
	bin, err := d.device.client.SPIR(d.device.handle, num)
	if err != nil {
		return 0, err
	}
	n := len(bin)
	copy(data, bin)
	return n, nil
}

func (d *pn532spiIO) MaxPacketSize() int {
	panic("TODO")
}
