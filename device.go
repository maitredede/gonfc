package gonfc

import (
	"time"

	"go.uber.org/zap"
)

type Device interface {
	ID() DeviceID
	Close() error

	SetPropertyBool(property Property, value bool) error
	SetPropertyInt(property Property, value int) error
	SetPropertyDuration(property Property, value time.Duration) error

	InitiatorInit() error
	// InitiatorInitSecureElement() error
	InitiatorSelectPassiveTarget(m Modulation, initData []byte) (*NfcTarget, error)
	InitiatorPollTarget(modulations []Modulation, pollnr byte, period byte) (*NfcTarget, error)
	// InitiatorSelectDepTarget(ndm DepMode, nbr BaudRate, pndiInitiator *DepInfo, timeout time.Duration) (*NfcTarget, error)
	InitiatorDeselectTarget() error
	InitiatorTransceiveBytes(tx, rx []byte, timeout time.Duration) (int, error)
	// InitiatorTransceiveBits()
	// InitiatorTransceiveBytesTimed()
	// InitiatorTransceiveBitsTimed()
	InitiatorTargetIsPresent(nt *NfcTarget) (bool, error)

	//WIP
	SetLastError(err error)
	GetInfiniteSelect() bool
	Logger() *zap.SugaredLogger
	DeviceGetSupportedModulation(mode Mode) ([]ModulationType, error)
	GetSupportedBaudRate(nmt ModulationType) ([]BaudRate, error)
	GetSupportedBaudRateTargetMode(nmt ModulationType) ([]BaudRate, error)
}

type NFCDeviceCommon struct {
	LastError      error
	InfiniteSelect bool
	Par            bool
	EasyFraming    bool
}
