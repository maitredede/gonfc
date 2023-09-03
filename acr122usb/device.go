package acr122usb

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/maitredede/gonfc"
	"github.com/maitredede/gonfc/acr122"
	"github.com/maitredede/gonfc/pn53x"
	"github.com/maitredede/gonfc/utils"
	"go.uber.org/zap"
)

type Acr122UsbSupportedDevice struct {
	VID  gousb.ID
	PID  gousb.ID
	Name string
}

var UsbSupportedDevices []Acr122UsbSupportedDevice = []Acr122UsbSupportedDevice{
	{VID: 0x072F, PID: 0x2200, Name: "ACS ACR122"},
	{VID: 0x072F, PID: 0x90CC, Name: "Touchatag"},
	{VID: 0x072F, PID: 0x2214, Name: "ACS ACR1222"},
}

type Acr122UsbDevice struct {
	usbdev *gousb.Device

	id     *acr122DeviceID
	device *gousb.Device
	logger *zap.SugaredLogger

	epInAdr    gousb.EndpointAddress
	epIn       *gousb.InEndpoint
	epInStream *gousb.ReadStream

	epOutAdr    gousb.EndpointAddress
	epOut       *gousb.OutEndpoint
	epOutStream *gousb.WriteStream

	maxPacketSize int
	io            *acr122io

	intf            *gousb.Interface
	done            func()
	name            string
	timerCorrection int

	chip *pn53x.Chip

	gonfc.NFCDeviceCommon
}

var (
	_ gonfc.Device        = (*Acr122UsbDevice)(nil)
	_ acr122.ACR122Device = (*Acr122UsbDevice)(nil)
)

func (pnd *Acr122UsbDevice) Close() error {
	errs := make([]error, 0)
	if pnd.epInStream != nil {
		errs = append(errs, pnd.epInStream.Close())
	}
	if pnd.epOutStream != nil {
		errs = append(errs, pnd.epOutStream.Close())
	}
	if pnd.done != nil {
		pnd.done()
	}
	if pnd.usbdev != nil {
		errs = append(errs, pnd.usbdev.Close())
	}
	return errors.Join(errs...)
}

func (pnd *Acr122UsbDevice) ID() gonfc.DeviceID {
	return pnd.id
}

func (d *Acr122UsbDevice) Logger() *zap.SugaredLogger {
	return d.logger
}

func (d *Acr122UsbDevice) String() string {
	return d.id.String()
}

func (d *Acr122UsbDevice) SetLastError(err error) {
	d.LastError = err
}
func (d *Acr122UsbDevice) GetInfiniteSelect() bool {
	return d.InfiniteSelect
}

func (pnd *Acr122UsbDevice) Name() string {
	return pnd.name
}

func (pnd *Acr122UsbDevice) usbGetUsbDeviceName() string {
	if vendor, vok := usbid.Vendors[pnd.id.desc.Vendor]; vok {
		if product, pok := vendor.Product[pnd.id.desc.Product]; pok {
			return fmt.Sprintf("%s / %s", vendor.Name, product.Name)
		}
	}
	return pnd.id.deviceInfo.Name
}

// usbinit (libnfc.acr122_usb_init)
func (pnd *Acr122UsbDevice) usbinit() error {
	// pnd.logger.Debugf("usbinit enter")
	// defer pnd.logger.Debugf("usbinit exit")

	abtRxBuf := make([]byte, 255+sizeOfCcidHeader)

	if err := pnd.chip.SetPropertyDuration(gonfc.NP_TIMEOUT_COMMAND, 1*time.Second); err != nil {
		return fmt.Errorf("set timeout failed: %w", err)
	}

	// pnd.logger.Debugf("  usbinit: ACR122 Get Firmware Version")
	// vBuff := make([]byte, 10)
	// if _, err := pnd.usbSendAPDU(0x00, 0x48, 0x00, nil, 0, vBuff); err != nil {
	// 	return fmt.Errorf("apdu send failed: %w", err)
	// }
	// pnd.logger.Debugf("fw version: %v", string(vBuff))

	//TODO : Bi-Color LED and Buzzer Control

	// Power On ICC
	// pnd.logger.Debugf("  usbinit: power on ICC")
	ccidFrame := []byte{PC_to_RDR_IccPowerOn, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}
	if _, err := pnd.usbBulkWrite(ccidFrame, 1000); err != nil {
		return fmt.Errorf("bulk write failed: %w", err)
	}
	if _, err := pnd.usbBulkRead(abtRxBuf, 1000); err != nil {
		return fmt.Errorf("bulk read failed: %w", err)
	}

	//ACR122 PICC Operating Parameters
	pnd.logger.Debugf("ACR122 PICC Operating Parameters")
	if _, err := pnd.usbSendAPDU(0x00, 0x51, 0x00, nil, 0, abtRxBuf); err != nil {
		return fmt.Errorf("apdu send failed: %w", err)
	}

	// if err := pnd.chip.Init(); err != nil {
	// 	return fmt.Errorf("pn53x init failed: %w", err)
	// }

	var err error
	for i := 0; i < 3; i++ {
		err = pnd.chip.Init()
		if err == nil {
			break
		}
		pnd.logger.Debugf("PN532 init failed (attempt %d/3) : %v", i+1, err)
	}
	if err != nil {
		return fmt.Errorf("pn53x init failed: %w", err)
	}

	return nil
}

func (pnd *Acr122UsbDevice) usbSendAPDU(ins byte, p1 byte, p2 byte, data []byte, le byte, out []byte) (int, error) {
	// pnd.logger.Debugf("  build frame apdu ins=%02x p1=%02x p2=%02x data len=%v LE=%v", ins, p1, p2, len(data), le)
	frame, err := buildFrameFromAPDU(ins, p1, p2, data, le)
	if err != nil {
		return 0, fmt.Errorf("usbSendAPDU frame build: %w", err)
	}
	if _, err := pnd.usbBulkWrite(frame, 1000); err != nil {
		return 0, fmt.Errorf("usbSendAPDU bulk write: %w", err)
	}
	n, err := pnd.usbBulkRead(out, 1000)
	if err != nil {
		return 0, fmt.Errorf("usbSendAPDU bulk read: %w", err)
	}
	return n, nil
}

func (pnd *Acr122UsbDevice) usbSend(data []byte, timeout time.Duration) (int, error) {
	// pnd.logger.Debugf("usbSend enter")
	// defer pnd.logger.Debugf("usbSend exit")

	// pnd.logger.Debugf("  build frame tama (data len=%v)", len(data))
	frame, err := buildFrameFromTama(data)
	if err != nil {
		pnd.LastError = gonfc.NFC_EINVARG
		return 0, pnd.LastError
	}
	n, err := pnd.usbBulkWrite(frame, timeout)
	if err != nil {
		pnd.LastError = err
		return n, pnd.LastError
	}
	return n, nil
}

const (
	USB_TIMEOUT_PER_PASS time.Duration = 200 * time.Millisecond
	USB_INFINITE_TIMEOUT time.Duration = 0
)

func (pnd *Acr122UsbDevice) usbReceive(data []byte, timeout time.Duration) (int, error) {
	// pnd.logger.Debugf("usbReceive enter")
	// defer pnd.logger.Debugf("usbReceive exit")
	// l := pnd.logger.Named("usbReceive")
	szDataLen := len(data)

	offset := 0

	abtRxBuf := make([]byte, 255+sizeOfCcidHeader)
	var usbTimeout time.Duration
	var remainingTime time.Duration = timeout
read:
	// l.Debugf("  read iteration")
	if timeout == USB_INFINITE_TIMEOUT {
		usbTimeout = USB_TIMEOUT_PER_PASS
	} else {
		remainingTime -= USB_TIMEOUT_PER_PASS
		if remainingTime <= 0 {
			pnd.LastError = gonfc.NFC_ETIMEOUT
			return 0, pnd.LastError
		} else {
			usbTimeout = min(remainingTime, USB_TIMEOUT_PER_PASS)
		}
	}

	n, err := pnd.usbBulkRead(abtRxBuf, usbTimeout)
	// l.Debugf("  usbBulkRead n=%v err=%v", n, err)
	attemptedResponse := RDR_to_PC_DataBlock
	if os.IsTimeout(err) {
		if pnd.AbortFlag {
			pnd.AbortFlag = false
			pnd.usbAck()
			pnd.LastError = gonfc.NFC_EOPABORTED
			return 0, pnd.LastError
		} else {
			goto read
		}
	}
	if err != nil || n < 10 {
		pnd.logger.Debugf("  s=%v d=%v", string(abtRxBuf[:n]), utils.ToHexString(abtRxBuf[:n]))
		pnd.logger.Errorf("invalid RDR_to_PC_DataBlock frame n=%v err=%v", n, err)
		// try to interrupt current device state
		pnd.usbAck()
		pnd.LastError = gonfc.BuildNFC_EIO(err)
		return 0, pnd.LastError
	}
	if abtRxBuf[offset] != attemptedResponse {
		pnd.logger.Errorf("Frame header mismatch (read)")
		pnd.LastError = gonfc.NFC_EIO
		return 0, pnd.LastError
	}
	offset++
	iLen := int(abtRxBuf[offset])
	offset++

	//status := abtRxBuf[7]
	error := abtRxBuf[8]
	if iLen == 0 && error == 0xFE { // ICC_MUTE; XXX check for more errors
		// Do not check status; my ACR122U seemingly has status=0 in this case,
		// even though the spec says it should have had bmCommandStatus=1
		// and bmICCStatus=1.
		pnd.logger.Debugf("command timed out")
		// log_put(LOG_GROUP, LOG_CATEGORY, NFC_LOG_PRIORITY_DEBUG, "%s", "Command timed out")
		// pnd.LastError = gonfc.NFC_ETIMEOUT
		// return 0, pnd.LastError
		goto read
	}

	rxB10 := abtRxBuf[10]
	rxB11 := abtRxBuf[11]
	if !(iLen > 1 && rxB10 == 0xd5) { // In case we didn't get an immediate answer:
		if iLen != 2 {
			pnd.logger.Errorf("Wrong reply")
			pnd.LastError = gonfc.NFC_EIO
			return 0, pnd.LastError
		}
		if rxB10 != SW1_More_Data_Available {
			if rxB10 == SW1_Warning_with_NV_changed && rxB11 == PN53x_Specific_Application_Level_Error_Code {
				pnd.logger.Errorf("PN532 has detected an error at the application level")
			} else if rxB10 == SW1_Warning_with_NV_changed && rxB11 == 0x00 {
				pnd.logger.Errorf("PN532 didn't reply")
			} else {
				pnd.logger.Errorf("Unexpected Status Word (SW1: %02x SW2: %02x)", rxB10, rxB11)
			}
			pnd.LastError = gonfc.NFC_EIO
			return 0, pnd.LastError
		}
		n, err := pnd.usbSendAPDU(APDU_GetAdditionnalData, 0x00, 0x00, nil, rxB11, abtRxBuf)
		if os.IsTimeout(err) {
			if pnd.AbortFlag {
				pnd.AbortFlag = false
				pnd.usbAck()
				pnd.LastError = gonfc.NFC_EOPABORTED
				return 0, pnd.LastError
			} else {
				goto read // FIXME May cause some trouble on Touchatag, right ?
			}
		}
		if err != nil || n < 10 {
			// try to interrupt current device state
			pnd.usbAck()
			pnd.LastError = gonfc.BuildNFC_EIO(err)
			return 0, pnd.LastError
		}
	}

	offset = 0
	if abtRxBuf[offset] != attemptedResponse {
		pnd.logger.Errorf("Frame header mismatch (last)")
		pnd.LastError = gonfc.NFC_EIO
		return 0, pnd.LastError
	}
	offset++

	// XXX In CCID specification, len is a 32-bits (dword), do we need to decode more than 1 byte ? (0-255 bytes for PN532 reply)
	iLen = int(abtRxBuf[offset])
	offset++
	if (abtRxBuf[offset] != 0x00) && (abtRxBuf[offset+1] != 0x00) && (abtRxBuf[offset+2] != 0x00) {
		pnd.logger.Errorf("Not implemented: only 1-byte length is supported, please report this bug with a full trace.")
		pnd.LastError = gonfc.NFC_EIO
		return 0, pnd.LastError
	}
	offset += 3

	if iLen < 4 {
		pnd.logger.Errorf("Too small reply")
		pnd.LastError = gonfc.NFC_EIO
		return 0, pnd.LastError
	}
	iLen -= 4 // We skip 2 bytes for PN532 direction byte (D5) and command byte (CMD+1), then 2 bytes for APDU status (90 00).
	if iLen > szDataLen {
		pnd.logger.Errorf("Unable to receive data: buffer too small. (szDataLen: %v, len: %v)", szDataLen, iLen)
		pnd.LastError = gonfc.NFC_EOVFLOW
		return 0, pnd.LastError
	}

	// Skip CCID remaining bytes
	offset += 2 // bSlot and bSeq are not used
	offset += 2 // bStatus and bError is partially checked
	offset += 1 // bRFU should be 0x00

	// TFI + PD0 (CC+1)
	if abtRxBuf[offset] != 0xD5 {
		pnd.logger.Errorf("TFI Mismatch")
		pnd.LastError = gonfc.NFC_EIO
		return 0, pnd.LastError
	}
	offset += 1

	if abtRxBuf[offset] != pnd.chip.LastCommandByte()+1 {
		pnd.logger.Errorf("Command Code verification failed")
		pnd.LastError = gonfc.NFC_EIO
		return 0, pnd.LastError
	}
	offset += 1

	copy(data, abtRxBuf[offset:offset+iLen])

	return iLen, nil
}

func (pnd *Acr122UsbDevice) usbAck() error {
	pnd.logger.Debug("usbAck enter")
	defer pnd.logger.Debug("usbAck exit")

	ackFrame := []byte{byte(pn53x.GetFirmwareVersion)} // We can't send a PN532's ACK frame, so we use a normal command to cancel current command
	pnd.logger.Debug("ACR122 Abort")
	pnd.logger.Debugf("  build frame tama (data len=%v)", len(ackFrame))
	frame, err := buildFrameFromTama(ackFrame)
	if err != nil {
		return err
	}
	_, err = pnd.usbBulkWrite(frame, 1000)
	if err != nil {
		return err
	}
	abtRxBuf := make([]byte, 255+sizeOfCcidHeader)
	_, err = pnd.usbBulkRead(abtRxBuf, 1000)
	return err
}

func (pnd *Acr122UsbDevice) SetPropertyBool(property gonfc.Property, value bool) error {
	pnd.logger.Debugf("  setPropertyBool %v: %v", gonfc.PropertyInfos[property].Name, value)
	return pnd.chip.SetPropertyBool(property, value)
}

func (pnd *Acr122UsbDevice) SetPropertyInt(property gonfc.Property, value int) error {
	pnd.logger.Debugf("  setPropertyInt %v: %v", gonfc.PropertyInfos[property].Name, value)
	return pnd.chip.SetPropertyInt(property, value)
}

func (pnd *Acr122UsbDevice) SetPropertyDuration(property gonfc.Property, value time.Duration) error {
	pnd.logger.Debugf("  SetPropertyDuration %v: %v", gonfc.PropertyInfos[property].Name, value)
	return pnd.chip.SetPropertyDuration(property, value)
}

func (pnd *Acr122UsbDevice) InitiatorTargetIsPresent(nt *gonfc.NfcTarget) (bool, error) {
	return pnd.chip.InitiatorTargetIsPresent(nt)
}
