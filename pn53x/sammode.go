package pn53x

import "github.com/maitredede/gonfc"

type samMode byte

const (
	samModeNormal      samMode = 0x01
	samModeVirtualCard samMode = 0x02
	samModeWiredCard   samMode = 0x03
	samModeDualCard    samMode = 0x04
)

func (pnd *Chip) samConfiguration(mode samMode, timeout int) error {
	abtCmd := []byte{byte(SAMConfiguration), byte(mode), 0x00, 0x00}
	szCmd := len(abtCmd)
	if pnd.chipType != PN532 {
		pnd.lastError.Set(gonfc.NFC_EDEVNOTSUPP)
		return pnd.lastError.Get()
	}
	switch mode {
	case samModeNormal: // Normal mode
		fallthrough
	case samModeWiredCard: // Wired card mode
		szCmd = 2
		break
	case samModeVirtualCard: // Virtual card mode
		fallthrough
	case samModeDualCard: // Dual card mode
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
