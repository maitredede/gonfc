package pn53x

import (
	"errors"
	"fmt"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/compat"
	"go.uber.org/zap"
)

var (
	pn53xAckFrame   = []byte{0x00, 0x00, 0xff, 0x00, 0xff, 0x00}
	pn53xNackFrame  = []byte{0x00, 0x00, 0xff, 0xff, 0x00, 0x00}
	pn53xErrorFrame = []byte{0x00, 0x00, 0xff, 0x01, 0xff, 0x7f, 0x81, 0x00}
)

const ( // PN53X_REG_CIU_BitFraming
	SYMBOL_START_SEND   byte = 0x80
	SYMBOL_RX_ALIGN     byte = 0x70
	SYMBOL_TX_LAST_BITS byte = 0x07
)

const ( //   PN53X_REG_CIU_RxMode
	SYMBOL_RX_CRC_ENABLE byte = 0x80
	SYMBOL_RX_SPEED      byte = 0x70
	SYMBOL_RX_NO_ERROR   byte = 0x08
	SYMBOL_RX_MULTIPLE   byte = 0x04
	// RX_FRAMING follow same scheme than TX_FRAMING
	SYMBOL_RX_FRAMING byte = 0x03
)

const (
	// Registers and symbols masks used to covers parts within a register
	//
	//	PN53X_REG_CIU_TxMode
	SYMBOL_TX_CRC_ENABLE byte = 0x80
	SYMBOL_TX_SPEED      byte = 0x70
	// TX_FRAMING bits explanation:
	//   00 : ISO/IEC 14443A/MIFARE and Passive Communication mode 106 kbit/s
	//   01 : Active Communication mode
	//   10 : FeliCa and Passive Communication mode at 212 kbit/s and 424 kbit/s
	//   11 : ISO/IEC 14443B
	SYMBOL_TX_FRAMING byte = 0x03
)

const ( //   PN53X_REG_CIU_TxAuto
	SYMBOL_FORCE_100_ASK byte = 0x40
	SYMBOL_AUTO_WAKE_UP  byte = 0x20
	SYMBOL_INITIAL_RF_ON byte = 0x04
)

const (
	//   PN53X_REG_CIU_Status2
	SYMBOL_MF_CRYPTO1_ON byte = 0x08

	//   PN53X_REG_CIU_ManualRCV
	SYMBOL_PARITY_DISABLE byte = 0x10
)

const (
	// PN53X Support Byte flags
	SUPPORT_ISO14443A byte = 0x01
	SUPPORT_ISO14443B byte = 0x02
	SUPPORT_ISO18092  byte = 0x04
)

const (
	/**
	 * Start bytes, packet length, length checksum, direction, packet checksum and postamble are overhead
	 */
	// The TFI is considered part of the overhead
	PN53x_NORMAL_FRAME__DATA_MAX_LEN   int = 254
	PN53x_NORMAL_FRAME__OVERHEAD       int = 8
	PN53x_EXTENDED_FRAME__DATA_MAX_LEN int = 264
	PN53x_EXTENDED_FRAME__OVERHEAD     int = 11
	PN53x_ACK_FRAME__LEN               int = 6
)

var (
	pn53xSupportedModulationAsTarget []gonfc.ModulationType = []gonfc.ModulationType{gonfc.NMT_ISO14443A, gonfc.NMT_FELICA, gonfc.NMT_DEP}
)

type Chip struct {
	io     IO
	logger *zap.SugaredLogger

	chipType     ChipType
	firmwareText string

	powerMode            PowerMode
	operatingMode        OperatingMode
	samMode              samMode
	lastStatusByte       byte
	wbTrigged            bool
	wbData               []byte
	wbMask               []byte
	timeoutCommand       int
	timeoutAtr           int
	timeoutCommunication int
	progressivefield     bool
	ui8Parameters        byte
	ui8TxBits            byte
	bCrc                 bool
	bPar                 compat.BoolFieldGetSet
	bEasyFraming         compat.BoolFieldGetSet
	bInfiniteSelect      compat.BoolFieldGetSet
	lastError            compat.ErrorFieldGetSet
	bAutoIso14443_4      bool
	btSupportByte        byte

	currentTarget                     any
	supported_modulation_as_initiator []gonfc.ModulationType
	supported_modulation_as_target    []gonfc.ModulationType

	lastCommand Command
}

func NewChip(io IO, logger *zap.SugaredLogger, bInfiniteSelect compat.BoolFieldGetSet, lastError compat.ErrorFieldGetSet, bPar compat.BoolFieldGetSet, bEasyFraming compat.BoolFieldGetSet) (*Chip, error) {
	if io == nil {
		panic(errors.New("io is nil"))
	}
	if logger == nil {
		panic(errors.New("logger is nil"))
	}
	chip := &Chip{
		logger: logger,
		// Keep I/O functions
		io: io,
		// Set type to generic (means unknown)
		chipType: PN53x,
		// Set power mode to normal, if your device starts in LowVBat (ie. PN532
		// UART) the driver layer have to correctly set it.
		powerMode: PowerModeNormal,
		// PN53x starts in initiator mode
		operatingMode: OperatingModeInitiator,
		// Clear last status byte
		lastStatusByte: 0x00,
		// Set current target to NULL
		currentTarget: nil,
		// Set current sam_mode to normal mode
		samMode: samModeNormal,
		// WriteBack cache is clean
		wbTrigged: false,
		wbMask:    make([]byte, PN53X_CACHE_REGISTER_SIZE),
		wbData:    make([]byte, PN53X_CACHE_REGISTER_SIZE),
		// Set default command timeout (350 ms)
		timeoutCommand: 350,

		// Set default ATR timeout (103 ms)
		timeoutAtr: 103,

		// Set default communication timeout (52 ms)
		timeoutCommunication: 52,

		supported_modulation_as_initiator: nil,

		supported_modulation_as_target: nil,

		// Set default progressive field flag
		progressivefield: false,

		//driver field accessors
		bInfiniteSelect: bInfiniteSelect,
		lastError:       lastError,
		bPar:            bPar,
		bEasyFraming:    bEasyFraming,
	}
	logger.Debug("NewChip")
	return chip, nil
}

func (pnd *Chip) LastCommandByte() byte {
	return byte(pnd.lastCommand)
}

func (pnd *Chip) Init() error {
	// pnd.logger.Debug("Init")
	if err := pnd.decodeFirmwareVersion(); err != nil {
		return err
	}

	if pnd.supported_modulation_as_initiator == nil {
		pnd.supported_modulation_as_initiator = make([]gonfc.ModulationType, 0)
		if (pnd.btSupportByte & SUPPORT_ISO14443A) != 0 {
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443A)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_FELICA)
		}
		if (pnd.btSupportByte & SUPPORT_ISO14443B) != 0 {
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443B)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443BI)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443B2SR)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443B2CT)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443BICLASS)
		}
		if pnd.chipType != PN531 {
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_JEWEL)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_BARCODE)
		}
		pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_DEP)
	}
	if pnd.supported_modulation_as_target == nil {
		pnd.supported_modulation_as_target = pn53xSupportedModulationAsTarget
	}

	// CRC handling should be enabled by default as declared in nfc_device_new
	// which is the case by default for pn53x, so nothing to do here
	// Parity handling should be enabled by default as declared in nfc_device_new
	// which is the case by default for pn53x, so nothing to do here

	// We can't read these parameters, so we set a default config by using the SetParameters wrapper
	// Note: pn53x_SetParameters() will save the sent value in pnd->ui8Parameters cache
	if err := pnd.SetParameters(PARAM_AUTO_ATR_RES | PARAM_AUTO_RATS); err != nil {
		return err
	}

	if err := pnd.resetSettings(); err != nil {
		return err
	}
	return nil
}

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

func (pnd *Chip) SetParameters(value Param) error {
	pnd.logger.Debugf("SetParameters")
	v := byte(value)
	abtCmd := []byte{byte(SetParameters), v}
	if _, err := pnd.transceive(abtCmd, nil, -1); err != nil {
		return err
	}
	pnd.ui8Parameters = v
	return nil
}

func (pnd *Chip) SetPropertyInt(property gonfc.Property, value int) error {
	// pnd.logger.Debugf("  setPropertyInt %v: %v", gonfc.PropertyInfos[property].Name, value)
	switch property {
	case gonfc.NP_TIMEOUT_COMMAND:
		pnd.timeoutCommand = value
		break
	case gonfc.NP_TIMEOUT_ATR:
		pnd.timeoutAtr = value
		return pnd.RFConfiguration__Various_timings(pn53x_int_to_timeout(pnd.timeoutAtr), pn53x_int_to_timeout(pnd.timeoutCommunication))
	case gonfc.NP_TIMEOUT_COM:
		pnd.timeoutCommunication = value
		return pnd.RFConfiguration__Various_timings(pn53x_int_to_timeout(pnd.timeoutAtr), pn53x_int_to_timeout(pnd.timeoutCommunication))
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

func (pnd *Chip) SetPropertyBool(property gonfc.Property, value bool) error {
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

func (pnd *Chip) SetParametersEnable(ui8Parameter Param, bEnable bool) error {
	// pnd.logger.Debugf("SetParametersEnable")
	var ui8Value byte
	if bEnable {
		ui8Value = byte(pnd.ui8Parameters | byte(ui8Parameter))
	} else {
		ui8Value = byte(pnd.ui8Parameters & ^(byte(ui8Parameter)))
	}
	if ui8Value != pnd.ui8Parameters {
		return pnd.SetParameters(Param(ui8Value))
	}
	return nil
}

func (pnd *Chip) resetSettings() error {
	// pnd.logger.Debugf("resetSettings")
	pnd.ui8TxBits = 0
	// Reset the ending transmission bits register, it is unknown what the last tranmission used there
	if err := pnd.writeRegisterMask(PN53X_REG_CIU_BitFraming, SYMBOL_TX_LAST_BITS, 0x00); err != nil {
		return err
	}
	// Make sure we reset the CRC and parity to chip handling.
	if err := pnd.SetPropertyBool(gonfc.NP_HANDLE_CRC, true); err != nil {
		return err
	}
	if err := pnd.SetPropertyBool(gonfc.NP_HANDLE_PARITY, true); err != nil {
		return err
	}
	// Activate "easy framing" feature by default
	if err := pnd.SetPropertyBool(gonfc.NP_EASY_FRAMING, true); err != nil {
		return err
	}
	// Deactivate the CRYPTO1 cipher, it may could cause problems when still active
	if err := pnd.SetPropertyBool(gonfc.NP_ACTIVATE_CRYPTO1, false); err != nil {
		return err
	}
	return nil
}

func (pnd *Chip) RFConfiguration__RF_field(enable bool) error {
	var val byte = 0x00
	if enable {
		val = 0x01
	}
	abtCmd := []byte{byte(RFConfiguration), byte(RFCI_FIELD), val}
	_, err := pnd.transceive(abtCmd, nil, -1)
	return err
}

func pn53x_int_to_timeout(ms int) byte {
	var res byte
	if ms > 0 {
		res = 0x10
		for i := 3280; i > 1; i /= 2 {
			if ms > i {
				break
			}
			res--
		}
	}
	return res
}

func (pnd *Chip) RFConfiguration__Various_timings(fATR_RES_Timeout byte, fRetryTimeout byte) error {
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

func (pnd *Chip) RFConfiguration__MaxRetries(MxRtyATR byte, MxRtyPSL byte, MxRtyPassiveActivation byte) error {
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
