package periphio

import (
	"time"

	"github.com/maitredede/gonfc/pn53x"
)

type periphioI2Cio struct {
	dev *PeriphioI2CDevice
}

var _ pn53x.IO = (*periphioI2Cio)(nil)

func (d *periphioI2Cio) Send(data []byte, timeout time.Duration) (int, error) {
	panic("TODO")
}

func (d *periphioI2Cio) Receive(data []byte, timeout time.Duration) (int, error) {
	panic("TODO")
}

func (d *periphioI2Cio) MaxPacketSize() int {
	panic("TODO")
}
