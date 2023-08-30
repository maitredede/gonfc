package pn53x

import (
	"time"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/compat"
	"go.uber.org/zap"
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

func (pnd *chipCommon) LastCommandByte() byte {
	return byte(pnd.lastCommand)
}
