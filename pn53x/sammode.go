package pn53x

import (
	"time"

	"github.com/maitredede/gonfc"
)

type SamMode byte

const (
	SamModeNormal      SamMode = 0x01
	SamModeVirtualCard SamMode = 0x02
	SamModeWiredCard   SamMode = 0x03
	SamModeDualCard    SamMode = 0x04
)

func (pnd *Chip) PN532SAMConfiguration(mode SamMode, timeout time.Duration) error {
	abtCmd := []byte{byte(SAMConfiguration), byte(mode), 0x00, 0x00}
	szCmd := len(abtCmd)
	if pnd.chipType != PN532 {
		pnd.lastError.Set(gonfc.NFC_EDEVNOTSUPP)
		return pnd.lastError.Get()
	}
	switch mode {
	case SamModeNormal: // Normal mode
		fallthrough
	case SamModeWiredCard: // Wired card mode
		szCmd = 2
		break
	case SamModeVirtualCard: // Virtual card mode
		fallthrough
	case SamModeDualCard: // Dual card mode
		// TODO Implement timeout handling
		szCmd = 3
		break
	default:
		pnd.lastError.Set(gonfc.NFC_EINVARG)
		return pnd.lastError.Get()
	}
	pnd.samMode = mode

	_, err := pnd.transceive(abtCmd[:szCmd], nil, timeout)
	return err
}
