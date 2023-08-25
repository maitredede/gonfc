package acr122usb

import (
	"encoding/binary"
	"io"
)

type apduHeader struct {
	bClass byte
	bIns   byte
	bP1    byte
	bP2    byte
	bLen   byte
}

const sizeOfApduHeader int = 5

func (b *apduHeader) Write(w io.Writer) error {
	e := binary.LittleEndian
	if err := binary.Write(w, e, b.bClass); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.bIns); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.bP1); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.bP2); err != nil {
		return err
	}
	if err := binary.Write(w, e, b.bLen); err != nil {
		return err
	}
	return nil
}

func (b *apduHeader) Read(r io.Reader) error {
	e := binary.BigEndian
	if err := binary.Read(r, e, &b.bClass); err != nil {
		return err
	}
	if err := binary.Read(r, e, &b.bIns); err != nil {
		return err
	}
	if err := binary.Read(r, e, &b.bP1); err != nil {
		return err
	}
	if err := binary.Read(r, e, &b.bP2); err != nil {
		return err
	}
	if err := binary.Read(r, e, &b.bLen); err != nil {
		return err
	}
	return nil
}
