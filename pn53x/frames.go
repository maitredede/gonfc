package pn53x

import (
	"bytes"

	"github.com/maitredede/gonfc"
)

var (
	AckFrame   = []byte{0x00, 0x00, 0xff, 0x00, 0xff, 0x00}
	NackFrame  = []byte{0x00, 0x00, 0xff, 0xff, 0x00, 0x00}
	ErrorFrame = []byte{0x00, 0x00, 0xff, 0x01, 0xff, 0x7f, 0x81, 0x00}
)

func (c *Chip) CheckAckFrame(frame []byte) error {
	if len(frame) >= len(AckFrame) {
		if bytes.Equal(AckFrame, frame[:len(AckFrame)]) {
			c.logger.Debug("PN53x ACKed")
			return nil
		}
	}
	c.lastError.Set(gonfc.NFC_EIO)
	c.logger.Error("Unexpected PN53x reply!")
	return c.lastError.Get()
}

func (c *Chip) CheckErrorFrame(frame []byte) error {
	if len(frame) >= len(ErrorFrame) {
		if bytes.Equal(ErrorFrame, frame[:len(ErrorFrame)]) {
			c.logger.Debug("PN53x sent an error frame")
			c.lastError.Set(gonfc.NFC_EIO)
			return c.lastError.Get()
		}
	}
	return nil
}
