package magic

import (
	"github.com/maitredede/gonfc"
	"go.uber.org/zap"
)

const (
	MAX_FRAME_LEN          int  = 264
	SAK_FLAG_ATS_SUPPORTED byte = 0x20
	CASCADE_BIT            byte = 0x04
)

// ISO14443A Anti-Collision Commands
var (
	abtReqa      = []byte{0x26}
	strangeWupa  = []byte{0x40}
	backdoorTest = []byte{0x43}
	abtSelectAll = []byte{0x93, 0x20}
	abtSelectTag = []byte{0x93, 0x70, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	abtRats      = []byte{0xe0, 0x50, 0x00, 0x00}
	abtHalt      = []byte{0x50, 0x00, 0x00, 0x00}
)

// IsMagicCard is an implementation of https://github.com/nfc-tools/libnfc/pull/461
func IsMagicCard(device gonfc.Device) (bool, error) {
	l := device.Logger().Named("magic")

	// Configure the CRC
	if err := device.SetPropertyBool(gonfc.NP_HANDLE_CRC, false); err != nil {
		return false, err
	}

	// Use raw send/receive methods
	if err := device.SetPropertyBool(gonfc.NP_EASY_FRAMING, false); err != nil {
		return false, err
	}

	// Disable 14443-4 autoswitching
	if err := device.SetPropertyBool(gonfc.NP_AUTO_ISO14443_4, false); err != nil {
		return false, err
	}

	if err := transmitBits(l, device, strangeWupa, 7); err != nil {
		l.Debugf("This is NOT a backdoored rewritable UID card")
		return false, nil
	}
	if err := transmitBytes(l, device, backdoorTest); err != nil {
		l.Debugf("This is backdoored rewritable UID card")
		return true, nil
	}
	l.Debugf("This is NOT a backdoored rewritable UID card")
	return false, nil
}

func transmitBits(l *zap.SugaredLogger, device gonfc.Device, pbtx []byte, szTxBits int) error {
	abtRx := make([]byte, 264)
	_ /*n*/, err := device.InitiatorTransceiveBits(pbtx, szTxBits, nil, abtRx, nil)
	if err != nil {
		return err
	}
	return nil
}

func transmitBytes(l *zap.SugaredLogger, device gonfc.Device, pbtx []byte) error {
	abtRx := make([]byte, 264)
	_ /*n*/, err := device.InitiatorTransceiveBytes(pbtx, abtRx, 0)
	if err != nil {
		return err
	}
	return nil
}
