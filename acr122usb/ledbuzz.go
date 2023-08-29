package acr122usb

import (
	"bytes"
	"fmt"
	"io"

	"github.com/maitredede/gonfc/acr122"
	"github.com/maitredede/gonfc/utils"
)

// See ACR122 manual: "Bi-Color LED and Buzzer Control" section

var getLedStateFrame []byte = []byte{
	0x6b,                                           // CCID
	0x09,                                           // lenght of frame
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // padding
	// frame:
	0xff,                   // Class
	0x00,                   // INS
	0x40,                   // P1: Get LED state command
	0x00,                   // P2: LED state control
	0x04,                   // Lc
	0x00, 0x00, 0x00, 0x00, // Blinking duration control
}

func (pnd *Acr122UsbDevice) ledBuzz(p2 byte, blink [4]byte) (acr122.LedState, error) {
	var state acr122.LedState
	ins := byte(0x00)
	p1 := byte(0x40)
	frame, err := buildFrameFromAPDU(ins, p1, p2, blink[:], 0)
	if err != nil {
		return state, err
	}
	if _, err := pnd.usbBulkWrite(frame, 1000); err != nil {
		return state, err
	}
	abtRxBuf := make([]byte, 255+sizeOfCcidHeader)
	n, err := pnd.usbBulkRead(abtRxBuf, 1000)
	if err != nil {
		return state, err
	}
	head, content, err := decodeAPDUResponse(abtRxBuf[:n])
	if err != nil {
		return state, err
	}
	if len(content) != 2 {
		return state, fmt.Errorf("invalid response length %v (expected: 2)", len(content))
	}

	sw1, sw2 := content[0], content[1]
	if sw1 == 0x90 {
		// The operation completed successfully.
		state.Red = (sw2 & (1 << 0)) != 0
		state.Green = (sw2 & (1 << 1)) != 0
		return state, nil
	}
	if sw1 == 0x63 {
		return state, fmt.Errorf("the operation failed")
	}

	pnd.logger.Warnf("  LED state unknown response: h=%+v d=%v", head, utils.ToHexString(content))
	return state, fmt.Errorf("the operation failed")
}

func (pnd *Acr122UsbDevice) GetLedState() (acr122.LedState, error) {
	pnd.logger.Debugf("ACR122 Get LED state")

	p2 := byte(0x00)
	blink := [4]byte{0x00, 0x00, 0x00, 0x00}

	return pnd.ledBuzz(p2, blink)
}

func decodeAPDUResponse(bin []byte) (ccidHeader, []byte, error) {
	buff := bytes.NewBuffer(bin)
	var header ccidHeader
	if err := header.Read(buff); err != nil {
		return header, nil, err
	}
	remaining, err := io.ReadAll(buff)
	if err != nil {
		return header, nil, err
	}
	return header, remaining, nil
}

func (pnd *Acr122UsbDevice) SetLed(state acr122.LedState) (acr122.LedState, error) {
	// ins := byte(0x00)
	// p1 := byte(0x40)
	p2 := byte(0x00)
	blink := [4]byte{0x00, 0x00, 0x00, 0x00}

	if state.Red {
		p2 |= 1 << 0 //red led on
		p2 |= 1 << 2 //red led state update
	} else {
		p2 &= ^byte(1 << 0) //red led off
		p2 |= 1 << 2        //red led state update
	}
	if state.Green {
		p2 |= 1 << 1 //red led on
		p2 |= 1 << 3 //red led state update
	} else {
		p2 &= ^byte(1 << 1) //red led off
		p2 |= 1 << 3        //red led state update
	}

	return pnd.ledBuzz(p2, blink)
}
