package acr122usb

import (
	"bytes"
	"testing"
)

func TestCcidHeaderWrite(t *testing.T) {

	data := ccidHeader{
		bMessageType:     0x01,
		dwLength:         0x02030405,
		bSlot:            0x06,
		bSeq:             0x07,
		bMessageSpecific: [3]byte{0x08, 0x09, 0x0A},
	}
	buff := &bytes.Buffer{}

	if err := data.Write(buff); err != nil {
		t.Fatal(err)
	}
	if buff.Len() != sizeOfCcidHeader {
		t.Fatal("length mismatch")
	}

	expected := make([]byte, sizeOfCcidHeader)
	for i := byte(1); int(i) <= len(expected); i++ {
		expected[i-1] = i
	}
	real := buff.Bytes()
	if !bytes.Equal(expected, real) {
		t.Fatal("content mismatch")
	}
}

func TestCcidHeaderRead(t *testing.T) {
	data := make([]byte, sizeOfCcidHeader)
	for i := byte(1); int(i) <= len(data); i++ {
		data[i-1] = i
	}
	buff := bytes.NewBuffer(data)

	var real ccidHeader
	if err := real.Read(buff); err != nil {
		t.Fatal(err)
	}

	if real.bMessageType != 0x01 {
		t.Fatal("unexpected value")
	}
	if real.dwLength != 0x02030405 {
		t.Fatal("unexpected value")
	}
	if real.bSlot != 0x06 {
		t.Fatal("unexpected value")
	}
	if real.bSeq != 0x07 {
		t.Fatal("unexpected value")
	}
	bin := []byte{0x08, 0x09, 0x0A}
	if !bytes.Equal(real.bMessageSpecific[:], bin) {
		t.Fatal("unexpected value")
	}
}
