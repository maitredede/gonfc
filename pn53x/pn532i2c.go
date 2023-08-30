package pn53x

import (
	"time"

	"github.com/maitredede/gonfc/compat"
	"go.uber.org/zap"
)

const (
	PN532_BUFFER_LEN int = PN53x_EXTENDED_FRAME__DATA_MAX_LEN + PN53x_EXTENDED_FRAME__OVERHEAD

	/*
	 * When sending lots of data, the pn532 occasionally fails to respond in time.
	 * Since it happens so rarely, lets try to fix it by re-sending the data. This
	 * define allows for fine tuning the number of retries.
	 */
	PN532_SEND_RETRIES int = 3
	/*
	 * Bus free time (in ms) between a STOP condition and START condition. See
	 * tBuf in the PN532 data sheet, section 12.25: Timing for the I2C interface,
	 * table 320. I2C timing specification, page 211, rev. 3.2 - 2007-12-07.
	 */
	PN532_BUS_FREE_TIME time.Duration = 5 * time.Millisecond
)

var (
	/* preamble and start bytes, see pn532-internal.h for details */
	PN53X_PREAMBLE_AND_START []byte = []byte{0x00, 0x00, 0xff}
)

type PN532I2CChip struct {
	chipCommon

	timerCorrection int
	abortFlag       compat.BoolFieldGetSet
}

var (
	PN53X_preamble_and_start     []byte = []byte{0x00, 0x00, 0xff}
	PN53X_PREAMBLE_AND_START_LEN int    = len(PN53X_preamble_and_start)
)

func NewPN532I2CChip(logger *zap.SugaredLogger, io IO, bInfiniteSelect compat.BoolFieldGetSet, lastError compat.ErrorFieldGetSet, bPar compat.BoolFieldGetSet, bEasyFraming compat.BoolFieldGetSet, abortFlag compat.BoolFieldGetSet) (*PN532I2CChip, error) {
	c := &PN532I2CChip{
		chipCommon: chipCommon{
			logger: logger,
			// Keep I/O functions
			io: io,

			// SAMConfiguration command if needed to wakeup the chip and pn53x_SAMConfiguration check if the chip is a PN532
			chipType: PN532,
			// This device starts in LowVBat mode
			powerMode: PowerModeLowVBat,

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
		abortFlag: abortFlag,
	}

	// empirical tuning
	c.timerCorrection = 48

	//driverData
	c.abortFlag.Set(false)

	// if err := c.CheckCommunication(); err != nil {
	// 	return nil, err
	// }
	return c, nil
}
