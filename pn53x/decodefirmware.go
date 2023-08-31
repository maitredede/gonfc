package pn53x

import (
	"fmt"

	"github.com/maitredede/gonfc"
)

func (pnd *Chip) decodeFirmwareVersion() error {
	// pnd.logger.Debug("decodeFirmwareVersion enter")
	// defer pnd.logger.Debug("decodeFirmwareVersion exit")

	abtCmd := []byte{byte(GetFirmwareVersion)}
	abtFw := make([]byte, 4)

	szFwLen, err := pnd.transceive(abtCmd, abtFw, -1)
	if err != nil {
		return fmt.Errorf("decodeFirmwareVersion transceive: %w", err)
	}
	if szFwLen == 2 {
		pnd.chipType = PN531
	} else if szFwLen == 4 {
		if abtFw[0] == 0x32 { // PN532 version IC
			pnd.chipType = PN532
		} else if abtFw[0] == 0x33 { // PN533 version IC
			if abtFw[1] == 0x01 { // Sony ROM code
				pnd.chipType = RCS360
			} else {
				pnd.chipType = PN533
			}
		} else {
			// Unknown version IC
			return gonfc.NFC_ENOTIMPL
		}
	} else {
		// Unknown chip
		return gonfc.NFC_ENOTIMPL
	}
	// Convert firmware info in text, PN531 gives 2 bytes info, but PN532 and PN533 gives 4
	switch pnd.chipType {
	case PN531:
		pnd.firmwareText = fmt.Sprintf("PN531 v%d.%d", abtFw[0], abtFw[1])
		pnd.btSupportByte = SUPPORT_ISO14443A | SUPPORT_ISO18092
	case PN532:
		pnd.firmwareText = fmt.Sprintf("PN532 v%d.%d", abtFw[1], abtFw[2])
		pnd.btSupportByte = abtFw[3]
	case PN533:
		fallthrough
	case RCS360:
		pnd.firmwareText = fmt.Sprintf("PN533 v%d.%d", abtFw[1], abtFw[2])
		pnd.btSupportByte = abtFw[3]
	case PN53x:
		// Could not happend
	}
	pnd.logger.Debugf("  firmware: %s", pnd.firmwareText)
	return nil
}
