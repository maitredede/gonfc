package acr122usb

import (
	"encoding/binary"
	"io"
)

type ccidHeader struct {
	bMessageType     byte
	dwLength         uint32
	bSlot            byte
	bSeq             byte
	bMessageSpecific [3]byte
}

const sizeOfCcidHeader int = 10

func (b *ccidHeader) Write(w io.Writer) error {
	e := binary.LittleEndian
	if err := binary.Write(w, e, b.bMessageType); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.dwLength); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.bSlot); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.bSeq); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.bMessageSpecific); err != nil {
		return err
	}
	return nil
}

func (b *ccidHeader) Read(r io.Reader) error {
	e := binary.BigEndian
	if err := binary.Read(r, e, &b.bMessageType); err != nil {
		return err
	}
	if err := binary.Read(r, e, &b.dwLength); err != nil {
		return err
	}
	if err := binary.Read(r, e, &b.bSlot); err != nil {
		return err
	}
	if err := binary.Read(r, e, &b.bSeq); err != nil {
		return err
	}
	if err := binary.Read(r, e, &b.bMessageSpecific); err != nil {
		return err
	}
	return nil
}
