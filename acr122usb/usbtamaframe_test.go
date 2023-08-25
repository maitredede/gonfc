package acr122usb

import (
	"bytes"
	"testing"
)

func TestFrameTemplateTama(t *testing.T) {
	frame := newUsbTamaFrame()

	buff := &bytes.Buffer{}
	if err := frame.Write(buff); err != nil {
		t.Fatal(err)
	}

	bin := buff.Bytes()
	for i := 0; i < len(usbFrameTemplate); i++ {
		if bin[i] != usbFrameTemplate[i] {
			t.Logf("tpl=%v", toHexString(usbFrameTemplate))
			t.Logf("bin=%v", toHexString(bin))
			t.Fatalf("value mismatch at index %v (len=%v)", i, len(usbFrameTemplate))
		}
	}
}
