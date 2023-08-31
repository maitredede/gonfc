package pn53x

import (
	"errors"

	"github.com/maitredede/gonfc/compat"
	"go.uber.org/zap"
)

func NewChip(io IO, logger *zap.SugaredLogger, bInfiniteSelect compat.BoolFieldGetSet, lastError compat.ErrorFieldGetSet, bPar compat.BoolFieldGetSet, bEasyFraming compat.BoolFieldGetSet, abortFlag compat.BoolFieldGetSet) (*Chip, error) {
	return newChipCommon(io, logger, bInfiniteSelect, lastError, bPar, bEasyFraming, abortFlag)
}

func NewChipPN532(io IO, logger *zap.SugaredLogger, bInfiniteSelect compat.BoolFieldGetSet, lastError compat.ErrorFieldGetSet, bPar compat.BoolFieldGetSet, bEasyFraming compat.BoolFieldGetSet, abortFlag compat.BoolFieldGetSet) (*Chip, error) {
	c, err := newChipCommon(io, logger, bInfiniteSelect, lastError, bPar, bEasyFraming, abortFlag)
	if err != nil {
		return nil, err
	}

	// SAMConfiguration command if needed to wakeup the chip and pn53x_SAMConfiguration check if the chip is a PN532
	c.chipType = PN532
	// This device starts in LowVBat mode
	c.powerMode = PowerModeLowVBat

	c.abortFlag.Set(false)

	return c, nil
}

func newChipCommon(io IO, logger *zap.SugaredLogger, bInfiniteSelect compat.BoolFieldGetSet, lastError compat.ErrorFieldGetSet, bPar compat.BoolFieldGetSet, bEasyFraming compat.BoolFieldGetSet, abortFlag compat.BoolFieldGetSet) (*Chip, error) {
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
		abortFlag:       abortFlag,

		supported_modulation_as_initiator: nil,
		supported_modulation_as_target:    nil,

		// PN53x starts in initiator mode
		operatingMode: OperatingModeInitiator,

		// Set current target to NULL
		currentTarget: nil,

		// Set default progressive field flag
		progressivefield: false,
		// empirical tuning
		timerCorrection: 48,
	}
	logger.Debug("NewChip")
	return chip, nil
}
