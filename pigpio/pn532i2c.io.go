package pigpio

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/pn53x"
)

type pn532i2cIO struct {
	device *PN532PiGPIOI2CDevice
}

var _ pn53x.IO = (*pn532i2cIO)(nil)

// Send
// drivers/pn532_i2c.c pn532_i2c_send
func (d *pn532i2cIO) Send(data []byte, timeout time.Duration) (int, error) {

	// Discard any existing data ?

	switch d.device.chip.PowerMode() {
	case pn53x.PowerModeLowVBat:
		/** PN532C106 wakeup. */
		if err := d.device.WakeUp(); err != nil {
			return 0, err
		}
		// According to PN532 application note, C106 appendix: to go out Low Vbat mode and enter in normal mode we need to send a SAMConfiguration command
		if err := d.device.chip.PN532SAMConfiguration(pn53x.SamModeNormal, 1000); err != nil {
			return 0, err
		}
	case pn53x.PowerModePowerDown:
		if err := d.device.WakeUp(); err != nil {
			return 0, err
		}
	case pn53x.PowerModeNormal:
		// Nothing to do :)
	}

	//abtFrame := make([]byte, pn53x.PN532_BUFFER_LEN)
	// // Every packet must start with the preamble and start bytes.
	// for i := 0; i < len(pn53x.PN53X_PREAMBLE_AND_START); i++ {
	// 	abtFrame[i] = pn53x.PN53X_PREAMBLE_AND_START[i]
	// }
	// if err := pn53x.BuildFrame(abtFrame, data); err != nil {
	// 	d.device.LastError = err
	// 	return 0, d.device.LastError
	// }
	abtFrame, err := pn53x.BuildFrame(data)
	if err != nil {
		return 0, err
	}

	var n int
	for retries := pn53x.PN532_SEND_RETRIES; retries > 0; retries-- {
		n, err = d.device.i2cWrite(abtFrame)
		if err != nil {
			d.device.logger.Errorf("Failed to transmit data. Retries left: %d.", retries-1)
			continue
		} else {
			break
		}
	}

	if err != nil {
		d.device.logger.Error("Unable to transmit data. (TX)")
		d.device.LastError = err
		return n, d.device.LastError
	}

	abtRxBuf := make([]byte, pn53x.PN53x_ACK_FRAME__LEN)

	// Wait for the ACK frame
	n, err = d.device.waitRdyFrame(abtRxBuf, timeout)
	if err != nil {
		if errors.Is(err, gonfc.NFC_EOPABORTED) {
			// Send an ACK frame from host to abort the command.
			d.device.i2cAck()
		}
		d.device.LastError = err
		return 0, d.device.LastError
	}

	if err := d.device.chip.CheckAckFrame(abtRxBuf[:n]); err != nil {
		return 0, err
	}
	return 0, nil
}

// Receive Read a response frame from the PN532 device
// drivers/pn532_i2c.c pn532_i2c_receive
func (d *pn532i2cIO) Receive(data []byte, timeout time.Duration) (int, error) {
	frameBuf := make([]byte, pn53x.PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	frameLength, err := d.device.waitRdyFrame(frameBuf, timeout)
	if errors.Is(err, gonfc.NFC_EOPABORTED) {
		return d.device.i2cAck()
	}
	if err != nil {
		d.device.LastError = err
		return 0, d.device.LastError
	}
	if !bytes.Equal(pn53x.PN53X_preamble_and_start, frameBuf[:pn53x.PN53X_PREAMBLE_AND_START_LEN]) {
		d.device.logger.Error("Frame preamble+start code mismatch")
		d.device.LastError = gonfc.NFC_EIO
		return 0, d.device.LastError
	}

	if frameBuf[3] == 0x01 && frameBuf[4] == 0xff {
		errorCode := frameBuf[5]
		err := fmt.Errorf("Application level error detected  (%d)", errorCode)
		d.device.logger.Error(err)
		d.device.LastError = gonfc.BuildNFC_EIO(err)
		return 0, d.device.LastError
	}

	var lg int
	var TFI_idx int
	if frameBuf[3] == 0xff && frameBuf[4] == 0xff {
		// Extended frame
		lg = int(frameBuf[5])<<8 + int(frameBuf[6])
		// Verify length checksum
		if ((int(frameBuf[5]) + int(frameBuf[6]) + int(frameBuf[7])) % 256) != 0 {
			err := errors.New("length checksum mismatch")
			d.device.logger.Error(err)
			d.device.LastError = gonfc.BuildNFC_EIO(err)
			return 0, d.device.LastError
		}
		TFI_idx = 8
	} else {
		// Normal frame
		lg = int(frameBuf[3])

		// Verify length checksum
		if (frameBuf[3] + frameBuf[4]) != 0 {
			err := errors.New("length checksum mismatch")
			d.device.logger.Error(err)
			d.device.LastError = gonfc.BuildNFC_EIO(err)
			return 0, d.device.LastError
		}
		TFI_idx = 5
	}

	if lg-2 > len(data) {
		err := fmt.Errorf("Unable to receive data: buffer too small. (szDataLen: %v, len: %v)", len(data), lg)
		d.device.logger.Error(err)
		d.device.LastError = gonfc.BuildNFC_EIO(err)
		return 0, d.device.LastError
	}

	TFI := frameBuf[TFI_idx]
	if TFI != 0xd5 {
		err := errors.New("TFI Mismatch")
		d.device.logger.Error(err)
		d.device.LastError = gonfc.BuildNFC_EIO(err)
		return 0, d.device.LastError
	}

	tfiNext := frameBuf[TFI_idx+1]
	lastCmdNext := d.device.chip.LastCommandByte() + 1
	if tfiNext != lastCmdNext {
		err := fmt.Errorf("Command Code verification failed.  (got 0x%02x,  expected 0x%02x)", tfiNext, lastCmdNext)
		d.device.logger.Error(err)
		d.device.LastError = gonfc.BuildNFC_EIO(err)
		return 0, d.device.LastError
	}

	DCS := frameBuf[TFI_idx+lg]
	btDCS := DCS

	// Compute data checksum
	for i := 0; i < lg; i++ {
		btDCS += frameBuf[TFI_idx+i]
	}

	if btDCS != 0 {
		err := fmt.Errorf("Data checksum mismatch  (DCS = 0x%02x, btDCS = 0x%02x)", DCS, btDCS)
		d.device.logger.Error(err)
		d.device.LastError = gonfc.BuildNFC_EIO(err)
		return 0, d.device.LastError
	}

	if frameBuf[TFI_idx+lg+1] != 0x00 {
		err := fmt.Errorf("Frame postamble mismatch  (got %d)", frameBuf[frameLength-1])
		d.device.logger.Error(err)
		d.device.LastError = gonfc.BuildNFC_EIO(err)
		return 0, d.device.LastError
	}

	for i := 0; i < lg-2; i++ {
		data[i] = frameBuf[i+TFI_idx+2]
	}
	return lg - 2, nil
}

func (d *pn532i2cIO) MaxPacketSize() int {
	panic("TODO")
}
