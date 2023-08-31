package pn53x

import (
	"bytes"

	"github.com/maitredede/gonfc"
)

func (c *Chip) CheckCommunication() error {
	abtCmd := []byte{byte(Diagnose), 0x00, 'g', 'o', '-', 'n', 'f', 'c'}
	abtExpectedRx := []byte{0x00, 'g', 'o', '-', 'n', 'f', 'c'}

	abtRx := make([]byte, len(abtExpectedRx))
	n, err := c.transceive(abtCmd, abtRx, 500)
	if err != nil {
		return err
	}

	received := abtRx[:n]
	if !bytes.Equal(received, abtExpectedRx) {
		return gonfc.NFC_EIO
	}
	return nil
}
