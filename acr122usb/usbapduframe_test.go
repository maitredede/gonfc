package acr122usb

import (
	"bytes"
	"testing"

	"github.com/maitredede/gonfc/utils"
)

func TestFrameTemplateAPDU(t *testing.T) {
	frame := newUsbApduFrame()

	buff := &bytes.Buffer{}
	if err := frame.Write(buff); err != nil {
		t.Fatal(err)
	}

	bin := buff.Bytes()
	for i := 0; i < len(usbFrameTemplate); i++ {
		if bin[i] != usbFrameTemplate[i] {
			t.Logf("tpl=%v", utils.ToHexString(usbFrameTemplate))
			t.Logf("bin=%v", utils.ToHexString(bin))
			t.Fatalf("value mismatch at index %v (len=%v)", i, len(usbFrameTemplate))
		}
	}
}
