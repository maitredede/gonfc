package pn53x

import (
	"time"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/compat"
	"go.uber.org/zap"
)

type Chip struct {
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
	timerCorrection      int
	ui8TxBits            byte
	ui8Parameters        byte
	bCrc                 bool
	samMode              SamMode
	lastCommand          Command
	lastStatusByte       byte
	bAutoIso14443_4      bool
	currentTarget        *gonfc.NfcTarget
	progressivefield     bool

	abortFlag       compat.BoolFieldGetSet
	bPar            compat.BoolFieldGetSet
	lastError       compat.ErrorFieldGetSet
	bEasyFraming    compat.BoolFieldGetSet
	bInfiniteSelect compat.BoolFieldGetSet
}

func (pnd *Chip) LastCommandByte() byte {
	return byte(pnd.lastCommand)
}

// currentTargetIs
// chips/pn53x.c pn53x_current_target_is
func (pnd *Chip) currentTargetIs(nt *gonfc.NfcTarget) bool {
	return pnd.currentTarget.Equals(nt)
}
