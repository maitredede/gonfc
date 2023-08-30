package pn53x

import (
	"errors"
	"time"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/compat"
	"go.uber.org/zap"
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

type chipCommon struct {
	io     IO
	logger *zap.SugaredLogger

	chipType  ChipType
	powerMode PowerMode

	firmwareText string

	btSupportByte                     byte
	supported_modulation_as_initiator []gonfc.ModulationType
	supported_modulation_as_target    []gonfc.ModulationType

	wbTrigged bool
	wbMask    []byte
	wbData    []byte

	timeoutCommand       time.Duration
	timeoutAtr           time.Duration
	timeoutCommunication time.Duration
	operatingMode        OperatingMode

	ui8TxBits     byte
	ui8Parameters byte
	bCrc          bool
	samMode       SamMode
	bPar          compat.BoolFieldGetSet

	lastCommand    Command
	lastStatusByte byte
	lastError      compat.ErrorFieldGetSet

	bEasyFraming    compat.BoolFieldGetSet
	bInfiniteSelect compat.BoolFieldGetSet
	bAutoIso14443_4 bool
}

type Chip struct {
	chipCommon

	progressivefield bool

	currentTarget any
}

func NewChip(io IO, logger *zap.SugaredLogger, bInfiniteSelect compat.BoolFieldGetSet, lastError compat.ErrorFieldGetSet, bPar compat.BoolFieldGetSet, bEasyFraming compat.BoolFieldGetSet) (*Chip, error) {
	if io == nil {
		panic(errors.New("io is nil"))
	}
	if logger == nil {
		panic(errors.New("logger is nil"))
	}
	chip := &Chip{
		chipCommon: chipCommon{
			logger: logger,
			// Keep I/O functions
			io: io,
			// Set type to generic (means unknown)
			chipType: PN53x,
			// Set power mode to normal, if your device starts in LowVBat (ie. PN532
			// UART) the driver layer have to correctly set it.
			powerMode: PowerModeNormal,

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

			// Clear last status byte
			lastStatusByte: 0x00,

			// Set current sam_mode to normal mode
			samMode: SamModeNormal,

			//driver field accessors
			lastError:       lastError,
			bPar:            bPar,
			bInfiniteSelect: bInfiniteSelect,
			bEasyFraming:    bEasyFraming,

			supported_modulation_as_initiator: nil,
			supported_modulation_as_target:    nil,

			// PN53x starts in initiator mode
			operatingMode: OperatingModeInitiator,
		},

		// Set current target to NULL
		currentTarget: nil,

		// Set default progressive field flag
		progressivefield: false,
	}
	logger.Debug("NewChip")
	return chip, nil
}

func (pnd *chipCommon) LastCommandByte() byte {
	return byte(pnd.lastCommand)
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

func pn53x_duration_to_timeout(t time.Duration) byte {
	ms := int(t.Milliseconds())
	return pn53x_int_to_timeout(ms)
}
