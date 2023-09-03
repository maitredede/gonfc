package pigpio

import (
	"errors"
	"time"

	"github.com/maitredede/go-pigpiod"
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/pn53x"
	"go.uber.org/zap"
)

type PN532PiGPIOI2CDevice struct {
	id     *PN532PiGPIOI2CDeviceID
	client *pigpiod.Conn
	handle uint32
	logger *zap.SugaredLogger

	chip *pn53x.Chip

	gonfc.NFCDeviceCommon

	transactionStop time.Time
}

var _ gonfc.Device = (*PN532PiGPIOI2CDevice)(nil)

func (d *PN532PiGPIOI2CDevice) ID() gonfc.DeviceID {
	return d.id
}

func (d *PN532PiGPIOI2CDevice) Logger() *zap.SugaredLogger {
	return d.logger
}

func (d *PN532PiGPIOI2CDevice) String() string {
	return d.id.String()
}

func (d *PN532PiGPIOI2CDevice) SetLastError(err error) {
	d.LastError = err
}

func (d *PN532PiGPIOI2CDevice) GetInfiniteSelect() bool {
	return d.InfiniteSelect
}

func (d *PN532PiGPIOI2CDevice) Close() error {
	errs := make([]error, 0)
	_, err := d.client.I2CC(d.handle)
	errs = append(errs, err)
	err = d.client.Close()
	errs = append(errs, err)
	return errors.Join(errs...)
}

func (pnd *PN532PiGPIOI2CDevice) SetPropertyBool(property gonfc.Property, value bool) error {
	pnd.logger.Debugf("  setPropertyBool %v: %v", gonfc.PropertyInfos[property].Name, value)
	return pnd.chip.SetPropertyBool(property, value)
}

func (pnd *PN532PiGPIOI2CDevice) SetPropertyInt(property gonfc.Property, value int) error {
	pnd.logger.Debugf("  setPropertyInt %v: %v", gonfc.PropertyInfos[property].Name, value)
	return pnd.chip.SetPropertyInt(property, value)
}

func (pnd *PN532PiGPIOI2CDevice) SetPropertyDuration(property gonfc.Property, value time.Duration) error {
	pnd.logger.Debugf("  SetPropertyDuration %v: %v", gonfc.PropertyInfos[property].Name, value)
	return pnd.chip.SetPropertyDuration(property, value)
}

func (d *PN532PiGPIOI2CDevice) InitiatorInit() error {
	return d.chip.InitiatorInit()
}

func (d *PN532PiGPIOI2CDevice) InitiatorSelectPassiveTarget(m gonfc.Modulation, initData []byte) (*gonfc.NfcTarget, error) {
	return d.chip.InitiatorSelectPassiveTarget(m, initData)
}

func (d *PN532PiGPIOI2CDevice) InitiatorTransceiveBytes(tx, rx []byte, timeout time.Duration) (int, error) {
	return d.chip.InitiatorTransceiveBytes(tx, rx, timeout)
}

func (pnd *PN532PiGPIOI2CDevice) InitiatorDeselectTarget() error {
	return pnd.chip.InitiatorDeselectTarget()
}

func (pnd *PN532PiGPIOI2CDevice) InitiatorPollTarget(modulations []gonfc.Modulation, pollnr byte, period byte) (*gonfc.NfcTarget, error) {
	return pnd.chip.InitiatorPollTarget(modulations, byte(pollnr), byte(period))
}

func (pnd *PN532PiGPIOI2CDevice) InitiatorTargetIsPresent(nt *gonfc.NfcTarget) (bool, error) {
	return pnd.chip.InitiatorTargetIsPresent(nt)
}

// WakeUp godoc
// pn532_i2c_wakeup
func (d *PN532PiGPIOI2CDevice) WakeUp() error {
	d.chip.SetPowerMode(pn53x.PowerModeNormal)
	return nil
}

func (pnd *PN532PiGPIOI2CDevice) waitBusFree() {
	if pnd.transactionStop.IsZero() {
		return
	}
	// elasped := time.Since(pnd.transactionStop)
	deadline := pnd.transactionStop.Add(pn53x.PN532_BUS_FREE_TIME)
	now := time.Now()
	if now.After(deadline) {
		return
	}
	toWait := deadline.Sub(now)
	time.Sleep(toWait)
}

func (pnd *PN532PiGPIOI2CDevice) endBusUse() {
	pnd.transactionStop = time.Now()
}

func (pnd *PN532PiGPIOI2CDevice) i2cWrite(buff []byte) (int, error) {
	pnd.waitBusFree()
	defer pnd.endBusUse()

	err := pnd.client.I2CWD(pnd.handle, buff)
	return len(buff), err
}

// drivers/pn532_i2c.c pn532_i2c_read
func (pnd *PN532PiGPIOI2CDevice) i2cRead(buff []byte) (int, error) {
	pnd.waitBusFree()
	defer pnd.endBusUse()

	return pnd.busI2cRead(buff)
}

// buses/i2c.c i2c_read
func (pnd *PN532PiGPIOI2CDevice) busI2cRead(buff []byte) (int, error) {
	data, err := pnd.client.I2cReadDevice(pnd.handle, len(buff))
	if err != nil {
		return 0, err
	}
	copy(buff, data)
	return len(data), nil
}

// waitRdyFrame Read data from the PN532 device until getting a frame with RDY bit set
// drivers/pn532_i2c.c pn532_i2c_wait_rdyframe
func (pnd *PN532PiGPIOI2CDevice) waitRdyFrame(pbtData []byte, timeout time.Duration) (int, error) {
	szDataLen := len(pbtData)
	done := false
	start := time.Now()
	//i2cRx := make([]byte, pn53x.PN53x_EXTENDED_FRAME__DATA_MAX_LEN+1)
	i2cRx := make([]byte, szDataLen+1)
	var resN int
	var resErr error
	for {
		recCount, err := pnd.i2cRead(i2cRx)
		if pnd.AbortFlag {
			// Reset abort flag
			pnd.AbortFlag = false
			pnd.logger.Debug("Wait for a READY frame has been aborted.")
			return 0, gonfc.NFC_EOPABORTED
		}
		if err != nil {
			done = true
			resN, resErr = 0, gonfc.NFC_EIO
		} else {
			rdy := i2cRx[0]
			if (rdy & 0x01) != 0 {
				done = true
				resN, resErr = recCount-1, nil
				copyLength := min(resN, szDataLen)
				for i := 0; i < copyLength; i++ {
					pbtData[i] = i2cRx[i+1]
				}
			} else {
				/* Not ready yet. Check for elapsed timeout. */
				if timeout > 0 && time.Since(start) > timeout {
					resN, resErr = 0, gonfc.NFC_ETIMEOUT
					done = true
					pnd.logger.Debug("timeout reached with no READY frame.")
				}
			}
		}

		//while (!done)
		if done {
			break
		}
	}
	return resN, resErr
}

func (pnd *PN532PiGPIOI2CDevice) i2cAck() (int, error) {
	return pnd.i2cWrite(pn53x.AckFrame)
}

func (pnd *PN532PiGPIOI2CDevice) DeviceGetSupportedModulation(mode gonfc.Mode) ([]gonfc.ModulationType, error) {
	return pnd.chip.GetSupportedModulation(mode)
}

func (pnd *PN532PiGPIOI2CDevice) GetSupportedBaudRate(mt gonfc.ModulationType) ([]gonfc.BaudRate, error) {
	return pnd.chip.GetSupportedBaudRate(gonfc.N_INITIATOR, mt)
}

func (pnd *PN532PiGPIOI2CDevice) GetSupportedBaudRateTargetMode(mt gonfc.ModulationType) ([]gonfc.BaudRate, error) {
	return pnd.chip.GetSupportedBaudRate(gonfc.N_TARGET, mt)
}

func (pnd *PN532PiGPIOI2CDevice) InitiatorTransceiveBits(tx []byte, txBits int, txPar []byte, rx []byte, rxPar []byte) (int, error) {
	return pnd.chip.InitiatorTransceiveBits(tx, txBits, txPar, rx, rxPar)
}
