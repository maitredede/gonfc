package pn53x

import (
	"errors"
	"fmt"

	"github.com/maitredede/gonfc"
)

func (pnd *Chip) cmdTrace(cmd byte) {
	info := pn53xCommands[Command(cmd)]
	pnd.logger.Debugf("  cmd:%s", info.name)
}

func (pnd *Chip) transceive(writeData []byte, readData []byte, timeout int) (int, error) {
	// pnd.logger.Debugf("transceive enter")
	// defer pnd.logger.Debugf("transceive exit")

	mi := false
	res := 0
	if pnd.wbTrigged {
		// pnd.logger.Debugf("triggering writeBackRegister")
		if err := pnd.writebackRegister(); err != nil {
			return 0, fmt.Errorf("writebackRegister: %w", err)
		}
	}

	pnd.cmdTrace(writeData[0])
	if timeout > 0 {
		pnd.logger.Debugf("Timeout value: %d", timeout)
	} else if timeout == 0 {
		pnd.logger.Debugf("No timeout")
	} else if timeout == -1 {
		timeout = pnd.timeoutCommand
	} else {
		pnd.logger.Errorf("Invalid timeout value: %d", timeout)
	}

	abtRx := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	szRx := len(abtRx)

	// Ensure a minimal receiving buffers is available
	if len(readData) == 0 {
		readData = make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	}

	// Call the send/receice callback functions of the current driver
	if _, err := pnd.io.Send(writeData, timeout); err != nil {
		return 0, fmt.Errorf("io.send error: %w", err)
	}

	// Command is sent, we store the command
	txCmd := Command(writeData[0])
	pnd.lastCommand = txCmd

	// Handle power mode for PN532
	if (pnd.chipType == PN532) && (TgInitAsTarget == txCmd) { // PN532 automatically goes into PowerDown mode when TgInitAsTarget command will be sent
		pnd.powerMode = PowerModePowerDown
	}

	var err error
	if res, err = pnd.io.Receive(readData, timeout); err != nil {
		return 0, fmt.Errorf("io.receive error: %w", err)
	}

	if (pnd.chipType == PN532) && (TgInitAsTarget == txCmd) { // PN532 automatically wakeup on external RF field
		pnd.powerMode = PowerModeNormal // When TgInitAsTarget reply that means an external RF have waken up the chip
	}

	switch txCmd {
	case PowerDown:
		fallthrough
	case InDataExchange:
		fallthrough
	case InCommunicateThru:
		fallthrough
	case InJumpForPSL:
		fallthrough
	case InPSL:
		fallthrough
	case InATR:
		fallthrough
	case InSelect:
		fallthrough
	case InJumpForDEP:
		fallthrough
	case TgGetData:
		fallthrough
	case TgGetInitiatorCommand:
		fallthrough
	case TgSetData:
		fallthrough
	case TgResponseToInitiator:
		fallthrough
	case TgSetGeneralBytes:
		fallthrough
	case TgSetMetaData:
		if (readData[0] & 0x80) != 0x00 {
			// NAD detected
			//abort()
			panic(errors.New("NAD detected"))
		}
		//      if (pbtRx[0] & 0x40) { abort(); } // MI detected
		mi = (readData[0] & 0x40) != 0x00
		pnd.lastStatusByte = readData[0] & 0x3f
		break
	case Diagnose:
		if writeData[1] == 0x06 { // Diagnose: Card presence detection
			pnd.lastStatusByte = readData[0] & 0x3f
		} else {
			pnd.lastStatusByte = 0
		}
		break
	case InDeselect:
		fallthrough
	case InRelease:
		if pnd.chipType == RCS360 {
			// Error code is in pbtRx[1] but we ignore error code anyway
			// because other PN53x chips always return 0 on those commands
			pnd.lastStatusByte = 0
			break
		}
		pnd.lastStatusByte = readData[0] & 0x3f
		break
	case ReadRegister:
		fallthrough
	case WriteRegister:
		if pnd.chipType == PN533 {
			// PN533 prepends its answer by the status byte
			pnd.lastStatusByte = readData[0] & 0x3f
		} else {
			pnd.lastStatusByte = 0
		}
		break
	default:
		pnd.lastStatusByte = 0
	}

	for mi {
		var res2 int
		var err2 error
		abtRx2 := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
		// Send empty command to card
		if /*res2*/ _, err2 = pnd.io.Send(writeData[:2], timeout); err2 != nil {
			return 0, err2
		}
		if res2, err2 = pnd.io.Receive(abtRx2, timeout); err2 != nil {
			return 0, err2
		}
		mi = (abtRx2[0] & 0x40) != 0x00
		if res+res2-1 > szRx {
			pnd.lastStatusByte = byte(ESMALLBUF)
			break
		}
		//memcpy(pbtRx+res, abtRx2+1, res2-1)
		for x := 0; x < res2; x++ {
			readData[res+x] = abtRx2[1+x]
		}
		// Copy last status byte
		readData[0] = abtRx2[0]
		res += res2 - 1
	}

	szRx = res

	var retRes int
	var retErr error
	chipErr := pn53xError(pnd.lastStatusByte)
	switch chipErr {
	case 0:
		retRes = szRx
		break
	case ETIMEOUT:
		fallthrough
	case ECRC:
		fallthrough
	case EPARITY:
		fallthrough
	case EBITCOUNT:
		fallthrough
	case EFRAMING:
		fallthrough
	case EBITCOLL:
		fallthrough
	case ERFPROTO:
		fallthrough
	case ERFTIMEOUT:
		fallthrough
	case EDEPUNKCMD:
		fallthrough
	case EDEPINVSTATE:
		fallthrough
	case ENAD:
		fallthrough
	case ENFCID3:
		fallthrough
	case EINVRXFRAM:
		fallthrough
	case EBCC:
		fallthrough
	case ECID:
		retErr = gonfc.NFC_ERFTRANS
		break
	case ESMALLBUF:
		fallthrough
	case EOVCURRENT:
		fallthrough
	case EBUFOVF:
		fallthrough
	case EOVHEAT:
		fallthrough
	case EINBUFOVF:
		retErr = gonfc.NFC_ECHIP
		break
	case EINVPARAM:
		fallthrough
	case EOPNOTALL:
		fallthrough
	case ECMD:
		fallthrough
	case ENSECNOTSUPP:
		retErr = gonfc.NFC_EINVARG
		break
	case ETGREL:
		fallthrough
	case ECDISCARDED:
		retErr = gonfc.NFC_ETGRELEASED
		// pn53x_current_target_free(pnd)
		break
	case EMFAUTH:
		// When a MIFARE Classic AUTH fails, the tag is automatically in HALT state
		retErr = gonfc.NFC_EMFCAUTHFAIL
		break
	default:
		retErr = gonfc.NFC_ECHIP
		break
	}

	if retErr != nil {
		pnd.lastError.Set(retErr)
		pnd.logger.Debugf("Chip error: \"%v\" (%02x), returned error: \"%s\" (%d))", chipErr, pnd.lastStatusByte, retErr, retRes)
	} else {
		pnd.lastError.Set(nil)
	}
	return retRes, retErr
}
