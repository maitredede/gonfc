package acr122usb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type usbTamaFrame struct {
	ccidHeader  ccidHeader
	apduHeader  apduHeader
	tamaHeader  byte
	tamaPayload [254]byte
}

func (b *usbTamaFrame) Write(w io.Writer) error {
	e := binary.LittleEndian
	if err := b.ccidHeader.Write(w); err != nil {
		return err
	}
	if err := b.apduHeader.Write(w); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.tamaHeader); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.tamaPayload); err != nil {
		return err
	}
	return nil
}

func newUsbTamaFrame() usbTamaFrame {
	frameTemplate := bytes.NewBuffer(usbFrameTemplate)
	var frame usbTamaFrame
	if err := frame.ReadTemplate(frameTemplate); err != nil {
		panic(err)
	}
	return frame
}

func (b *usbTamaFrame) ReadTemplate(r io.Reader) error {
	if err := b.ccidHeader.Read(r); err != nil {
		return err
	}
	if err := b.apduHeader.Read(r); err != nil {
		return err
	}
	b.tamaHeader = 0
	for i := 0; i < len(b.tamaPayload); i++ {
		b.tamaPayload[i] = 0
	}

	br := make([]byte, 1)
	if _, err := r.Read(br); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	b.tamaHeader = br[0]
	for i := 0; i < len(b.tamaPayload); i++ {
		if _, err := r.Read(br); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		b.tamaPayload[i] = br[i]
	}
	return nil
}

func buildFrameFromTama(tama []byte) ([]byte, error) {
	frame := newUsbTamaFrame()
	tamaLen := len(tama)
	if tamaLen > len(frame.tamaPayload) {
		return nil, fmt.Errorf("data too long (len=%d max=%d)", tamaLen, len(frame.tamaPayload))
	}
	frame.ccidHeader.dwLength = uint32(tamaLen + sizeOfApduHeader + 1)
	frame.apduHeader.bLen = byte(tamaLen + 1)
	for i := 0; i < tamaLen; i++ {
		frame.tamaPayload[i] = tama[i]
	}
	out := &bytes.Buffer{}
	if err := frame.Write(out); err != nil {
		return nil, err
	}

	bin := out.Bytes()
	size := sizeOfCcidHeader + sizeOfApduHeader + 1 + tamaLen
	return bin[:size], nil
}
