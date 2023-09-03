package periphio

import (
	"time"

	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/pn53x"
	"go.uber.org/zap"
	"periph.io/x/conn/v3/i2c"
)

type PeriphioI2CDevice struct {
	id     *PeriphioI2CDeviceID
	logger *zap.SugaredLogger
	driver *PeriphioI2CDriver

	hwBus i2c.BusCloser
	hwDev *i2c.Dev
	io    *periphioI2Cio
	chip  *pn53x.Chip

	gonfc.NFCDeviceCommon
}

var _ gonfc.Device = (*PeriphioI2CDevice)(nil)

func (d *PeriphioI2CDevice) ID() gonfc.DeviceID {
	return d.id
}

func (d *PeriphioI2CDevice) Logger() *zap.SugaredLogger {
	return d.logger
}

func (d *PeriphioI2CDevice) Close() error {
	return d.hwBus.Close()
}

func (d *PeriphioI2CDevice) String() string {
	return d.id.String()
}

func (d *PeriphioI2CDevice) SetPropertyBool(property gonfc.Property, value bool) error {
	panic("TODO")
}

func (d *PeriphioI2CDevice) SetPropertyInt(property gonfc.Property, value int) error {
	panic("TODO")
}

func (d *PeriphioI2CDevice) SetPropertyDuration(property gonfc.Property, value time.Duration) error {
	panic("TODO")
}

func (d *PeriphioI2CDevice) InitiatorInit() error {
	panic("TODO")
}

func (d *PeriphioI2CDevice) InitiatorSelectPassiveTarget(m gonfc.Modulation, initData []byte) (*gonfc.NfcTarget, error) {
	panic("TODO")
}

func (d *PeriphioI2CDevice) InitiatorPollTarget(modulations []gonfc.Modulation, pollnr byte, period byte) (*gonfc.NfcTarget, error) {
	panic("TODO")
}

func (d *PeriphioI2CDevice) InitiatorDeselectTarget() error {
	panic("TODO")
}

func (d *PeriphioI2CDevice) InitiatorTransceiveBytes(tx, rx []byte, timeout time.Duration) (int, error) {
	panic("TODO")
}

func (d *PeriphioI2CDevice) InitiatorTargetIsPresent(nt *gonfc.NfcTarget) (bool, error) {
	panic("TODO")
}

// WIP
func (d *PeriphioI2CDevice) SetLastError(err error) {
	d.LastError = err
}

func (d *PeriphioI2CDevice) GetInfiniteSelect() bool {
	return d.InfiniteSelect
}

func (d *PeriphioI2CDevice) DeviceGetSupportedModulation(mode gonfc.Mode) ([]gonfc.ModulationType, error) {
	panic("TODO")
}

func (d *PeriphioI2CDevice) GetSupportedBaudRate(nmt gonfc.ModulationType) ([]gonfc.BaudRate, error) {
	panic("TODO")
}

func (d *PeriphioI2CDevice) GetSupportedBaudRateTargetMode(nmt gonfc.ModulationType) ([]gonfc.BaudRate, error) {
	panic("TODO")
}

func (pnd *PeriphioI2CDevice) InitiatorTransceiveBits(tx []byte, txBits int, txPar []byte, rx []byte, rxPar []byte) (int, error) {
	return pnd.chip.InitiatorTransceiveBits(tx, txBits, txPar, rx, rxPar)
}
