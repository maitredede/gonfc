package acr122usb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/maitredede/gonfc"
)

type usbApduFrame struct {
	ccidHeader  ccidHeader
	apduHeader  apduHeader
	apduPayload [255]byte
}

func (b *usbApduFrame) Write(w io.Writer) error {
	e := binary.LittleEndian
	if err := b.ccidHeader.Write(w); err != nil {
		return err
	}
	if err := b.apduHeader.Write(w); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.apduPayload); err != nil {
		return err
	}
	return nil
}

func (b *usbApduFrame) ReadTemplate(r io.Reader) error {
	if err := b.ccidHeader.Read(r); err != nil {
		return err
	}
	if err := b.apduHeader.Read(r); err != nil {
		return err
	}
	for i := 0; i < len(b.apduPayload); i++ {
		b.apduPayload[i] = 0
	}

	br := make([]byte, 1)
	for i := 0; i < len(b.apduPayload); i++ {
		if _, err := r.Read(br); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		b.apduPayload[i] = br[i]
	}
	return nil
}

func newUsbApduFrame() usbApduFrame {
	frameTemplate := bytes.NewBuffer(usbFrameTemplate)
	var frame usbApduFrame
	if err := frame.ReadTemplate(frameTemplate); err != nil {
		panic(err)
	}
	return frame
}

func buildFrameFromAPDU(ins, p1, p2 byte, data []byte, le byte) ([]byte, error) {
	dataLen := len(data)
	frame := newUsbApduFrame()
	if dataLen > len(frame.apduPayload) {
		return nil, gonfc.NFC_EINVARG
	}

	frame.ccidHeader.dwLength = uint32(dataLen + sizeOfApduHeader)
	frame.apduHeader.bIns = ins
	frame.apduHeader.bP1 = p1
	frame.apduHeader.bP2 = p2
	if dataLen > 0 {
		frame.apduHeader.bLen = byte(dataLen)
		for i := 0; i < dataLen; i++ {
			frame.apduPayload[i] = data[i]
		}
	} else {
		frame.apduHeader.bLen = le
	}

	out := &bytes.Buffer{}
	if err := frame.Write(out); err != nil {
		return nil, fmt.Errorf("frame write: %w", err)
	}
	bin := out.Bytes()
	size := sizeOfCcidHeader + sizeOfApduHeader + dataLen
	return bin[:size], nil
}
