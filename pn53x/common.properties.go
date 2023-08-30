package pn53x

import (
	"time"

	"github.com/maitredede/gonfc"
)

func (pnd *chipCommon) SetPropertyDuration(property gonfc.Property, value time.Duration) error {
	switch property {
	case gonfc.NP_TIMEOUT_COMMAND:
		pnd.timeoutCommand = value
		break
	case gonfc.NP_TIMEOUT_ATR:
		pnd.timeoutAtr = value
		return pnd.RFConfiguration__Various_timings(pn53x_duration_to_timeout(pnd.timeoutAtr), pn53x_duration_to_timeout(pnd.timeoutCommunication))
	case gonfc.NP_TIMEOUT_COM:
		pnd.timeoutCommunication = value
		return pnd.RFConfiguration__Various_timings(pn53x_duration_to_timeout(pnd.timeoutAtr), pn53x_duration_to_timeout(pnd.timeoutCommunication))
	}
	return gonfc.NFC_EINVARG
}

func (pnd *chipCommon) SetPropertyInt(property gonfc.Property, value int) error {
	// pnd.logger.Debugf("  setPropertyInt %v: %v", gonfc.PropertyInfos[property].Name, value)
	switch property {
	case gonfc.NP_TIMEOUT_COMMAND:
		//pnd.timeoutCommand = value
		fallthrough
	case gonfc.NP_TIMEOUT_ATR:
		// pnd.timeoutAtr = value
		// return pnd.RFConfiguration__Various_timings(pn53x_int_to_timeout(pnd.timeoutAtr), pn53x_int_to_timeout(pnd.timeoutCommunication))
		fallthrough
	case gonfc.NP_TIMEOUT_COM:
		// pnd.timeoutCommunication = value
		// return pnd.RFConfiguration__Various_timings(pn53x_int_to_timeout(pnd.timeoutAtr), pn53x_int_to_timeout(pnd.timeoutCommunication))
		pnd.logger.Debugf("SetPropertyInt called for %v, prefer SetPropertyDuration", gonfc.PropertyInfos[property].Name)
		return pnd.SetPropertyDuration(property, time.Duration(value)*time.Millisecond)
	// Following properties are invalid (not integer)
	case gonfc.NP_HANDLE_CRC:
		fallthrough
	case gonfc.NP_HANDLE_PARITY:
		fallthrough
	case gonfc.NP_ACTIVATE_FIELD:
		fallthrough
	case gonfc.NP_ACTIVATE_CRYPTO1:
		fallthrough
	case gonfc.NP_INFINITE_SELECT:
		fallthrough
	case gonfc.NP_ACCEPT_INVALID_FRAMES:
		fallthrough
	case gonfc.NP_ACCEPT_MULTIPLE_FRAMES:
		fallthrough
	case gonfc.NP_AUTO_ISO14443_4:
		fallthrough
	case gonfc.NP_EASY_FRAMING:
		fallthrough
	case gonfc.NP_FORCE_ISO14443_A:
		fallthrough
	case gonfc.NP_FORCE_ISO14443_B:
		fallthrough
	case gonfc.NP_FORCE_SPEED_106:
		return gonfc.NFC_EINVARG
	}
	return nil
}

func (pnd *chipCommon) SetPropertyBool(property gonfc.Property, value bool) error {
	// pnd.logger.Debugf("  setPropertyBool %v: %v", gonfc.PropertyInfos[property].Name, value)
	switch property {
	case gonfc.NP_HANDLE_CRC:
		if pnd.bCrc == value {
			return nil
		}
		// TX and RX are both represented by the symbol 0x80
		var btValue byte = 0x00
		if value {
			btValue = 0x80
		}
		if err := pnd.writeRegisterMask(PN53X_REG_CIU_TxMode, SYMBOL_TX_CRC_ENABLE, btValue); err != nil {
			return err
		}
		if err := pnd.writeRegisterMask(PN53X_REG_CIU_RxMode, SYMBOL_RX_CRC_ENABLE, btValue); err != nil {
			return err
		}
		pnd.bCrc = value
		return nil
	case gonfc.NP_HANDLE_PARITY:
		// Handle parity bit by PN53X chip or parse it as data bit
		if pnd.bPar.Get() == value {
			return nil
		}
		var btValue byte
		if value {
			btValue = 0x00
		} else {
			btValue = SYMBOL_PARITY_DISABLE
		}
		if err := pnd.writeRegisterMask(PN53X_REG_CIU_ManualRCV, SYMBOL_PARITY_DISABLE, btValue); err != nil {
			return err
		}
		pnd.bPar.Set(value)
		return nil
	case gonfc.NP_EASY_FRAMING:
		pnd.bEasyFraming.Set(value)
		return nil

	case gonfc.NP_ACTIVATE_FIELD:
		return pnd.RFConfiguration__RF_field(value)

	case gonfc.NP_ACTIVATE_CRYPTO1:
		var btValue byte = 0x00
		if value {
			btValue = SYMBOL_MF_CRYPTO1_ON
		}
		return pnd.writeRegisterMask(PN53X_REG_CIU_Status2, SYMBOL_MF_CRYPTO1_ON, btValue)

	case gonfc.NP_INFINITE_SELECT:
		// TODO Made some research around this point:
		// timings could be tweak better than this, and maybe we can tweak timings
		// to "gain" a sort-of hardware polling (ie. like PN532 does)
		pnd.bInfiniteSelect.Set(value)
		var valMxRtyATR byte = 0xff
		var valMxRtyPSL byte = 0xff
		var valMxRtyPassiveActivation byte = 0xff
		if value {
			valMxRtyATR = 0x00
			valMxRtyPSL = 0x01
			valMxRtyPassiveActivation = 0x02
		}
		return pnd.RFConfiguration__MaxRetries(valMxRtyATR, valMxRtyPSL, valMxRtyPassiveActivation)
	case gonfc.NP_ACCEPT_INVALID_FRAMES:
		var btValue byte = 0x00
		if value {
			btValue = SYMBOL_RX_NO_ERROR
		}
		return pnd.writeRegisterMask(PN53X_REG_CIU_RxMode, SYMBOL_RX_NO_ERROR, btValue)

	case gonfc.NP_ACCEPT_MULTIPLE_FRAMES:
		var btValue byte = 0x00
		if value {
			btValue = SYMBOL_RX_MULTIPLE
		}
		return pnd.writeRegisterMask(PN53X_REG_CIU_RxMode, SYMBOL_RX_MULTIPLE, btValue)

	case gonfc.NP_AUTO_ISO14443_4:
		if value == pnd.bAutoIso14443_4 {
			// Nothing to do
			return nil
		}
		pnd.bAutoIso14443_4 = value
		return pnd.SetParametersEnable(PARAM_AUTO_RATS, value)

	case gonfc.NP_FORCE_ISO14443_A:
		if !value {
			// Nothing to do
			return nil
		}
		// Force pn53x to be in ISO14443-A mode
		if err := pnd.writeRegisterMask(PN53X_REG_CIU_TxMode, SYMBOL_TX_FRAMING, 0x00); err != nil {
			return err
		}
		if err := pnd.writeRegisterMask(PN53X_REG_CIU_RxMode, SYMBOL_RX_FRAMING, 0x00); err != nil {
			return err
		}

		// Set the PN53X to force 100% ASK Modified miller decoding (default for 14443A cards)
		return pnd.writeRegisterMask(PN53X_REG_CIU_TxAuto, SYMBOL_FORCE_100_ASK, 0x40)

	case gonfc.NP_FORCE_ISO14443_B:
		if !value {
			// Nothing to do
			return nil
		}
		// Force pn53x to be in ISO14443-B mode
		if err := pnd.writeRegisterMask(PN53X_REG_CIU_TxMode, SYMBOL_TX_FRAMING, 0x03); err != nil {
			return err
		}
		return pnd.writeRegisterMask(PN53X_REG_CIU_RxMode, SYMBOL_RX_FRAMING, 0x03)

	case gonfc.NP_FORCE_SPEED_106:
		if !value {
			// Nothing to do
			return nil
		}
		// Force pn53x to be at 106 kbps
		if err := pnd.writeRegisterMask(PN53X_REG_CIU_TxMode, SYMBOL_TX_SPEED, 0x00); err != nil {
			return err
		}
		return pnd.writeRegisterMask(PN53X_REG_CIU_RxMode, SYMBOL_RX_SPEED, 0x00)
	// Following properties are invalid (not boolean)
	case gonfc.NP_TIMEOUT_COMMAND:
		fallthrough
	case gonfc.NP_TIMEOUT_ATR:
		fallthrough
	case gonfc.NP_TIMEOUT_COM:
		return gonfc.NFC_EINVARG
	}

	return gonfc.NFC_EINVARG
}

func (pnd *chipCommon) RFConfiguration__Various_timings(fATR_RES_Timeout byte, fRetryTimeout byte) error {
	abtCmd := []byte{
		byte(RFConfiguration),
		byte(RFCI_TIMING),
		0x00, //RFU
		fATR_RES_Timeout,
		fRetryTimeout,
	}
	_, err := pnd.transceive(abtCmd, nil, -1)
	return err
}

func (pnd *chipCommon) RFConfiguration__MaxRetries(MxRtyATR byte, MxRtyPSL byte, MxRtyPassiveActivation byte) error {
	abtCmd := []byte{
		byte(RFConfiguration),
		byte(RFCI_RETRY_SELECT),
		MxRtyATR,
		MxRtyPSL,
		MxRtyPassiveActivation,
	}
	_, err := pnd.transceive(abtCmd, nil, -1)
	return err
}

func (pnd *chipCommon) RFConfiguration__RF_field(enable bool) error {
	var val byte = 0x00
	if enable {
		val = 0x01
	}
	abtCmd := []byte{byte(RFConfiguration), byte(RFCI_FIELD), val}
	_, err := pnd.transceive(abtCmd, nil, -1)
	return err
}
