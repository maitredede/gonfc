package pn53x

import "github.com/maitredede/gonfc"

const (
	PN53X_CACHE_REGISTER_MIN_ADDRESS = uint16(PN53X_REG_CIU_Mode)
	PN53X_CACHE_REGISTER_MAX_ADDRESS = uint16(PN53X_REG_CIU_Coll)
	PN53X_CACHE_REGISTER_SIZE        = int((PN53X_CACHE_REGISTER_MAX_ADDRESS - PN53X_CACHE_REGISTER_MIN_ADDRESS) + 1)
)

type Register uint16

const (
	// Register addresses
	PN53X_REG_Control_switch_rng Register = 0x6106
	PN53X_REG_CIU_Mode           Register = 0x6301
	PN53X_REG_CIU_TxMode         Register = 0x6302
	PN53X_REG_CIU_RxMode         Register = 0x6303
	PN53X_REG_CIU_TxControl      Register = 0x6304
	PN53X_REG_CIU_TxAuto         Register = 0x6305
	PN53X_REG_CIU_TxSel          Register = 0x6306
	PN53X_REG_CIU_RxSel          Register = 0x6307
	PN53X_REG_CIU_RxThreshold    Register = 0x6308
	PN53X_REG_CIU_Demod          Register = 0x6309
	PN53X_REG_CIU_FelNFC1        Register = 0x630A
	PN53X_REG_CIU_FelNFC2        Register = 0x630B
	PN53X_REG_CIU_MifNFC         Register = 0x630C
	PN53X_REG_CIU_ManualRCV      Register = 0x630D
	PN53X_REG_CIU_TypeB          Register = 0x630E
	// #define PN53X_REG_- 0x630F
	// #define PN53X_REG_- 0x6310
	PN53X_REG_CIU_CRCResultMSB   Register = 0x6311
	PN53X_REG_CIU_CRCResultLSB   Register = 0x6312
	PN53X_REG_CIU_GsNOFF         Register = 0x6313
	PN53X_REG_CIU_ModWidth       Register = 0x6314
	PN53X_REG_CIU_TxBitPhase     Register = 0x6315
	PN53X_REG_CIU_RFCfg          Register = 0x6316
	PN53X_REG_CIU_GsNOn          Register = 0x6317
	PN53X_REG_CIU_CWGsP          Register = 0x6318
	PN53X_REG_CIU_ModGsP         Register = 0x6319
	PN53X_REG_CIU_TMode          Register = 0x631A
	PN53X_REG_CIU_TPrescaler     Register = 0x631B
	PN53X_REG_CIU_TReloadVal_hi  Register = 0x631C
	PN53X_REG_CIU_TReloadVal_lo  Register = 0x631D
	PN53X_REG_CIU_TCounterVal_hi Register = 0x631E
	PN53X_REG_CIU_TCounterVal_lo Register = 0x631F
	// #define PN53X_REG_- 0x6320
	PN53X_REG_CIU_TestSel1     Register = 0x6321
	PN53X_REG_CIU_TestSel2     Register = 0x6322
	PN53X_REG_CIU_TestPinEn    Register = 0x6323
	PN53X_REG_CIU_TestPinValue Register = 0x6324
	PN53X_REG_CIU_TestBus      Register = 0x6325
	PN53X_REG_CIU_AutoTest     Register = 0x6326
	PN53X_REG_CIU_Version      Register = 0x6327
	PN53X_REG_CIU_AnalogTest   Register = 0x6328
	PN53X_REG_CIU_TestDAC1     Register = 0x6329
	PN53X_REG_CIU_TestDAC2     Register = 0x632A
	PN53X_REG_CIU_TestADC      Register = 0x632B
	// #define PN53X_REG_- 0x632C
	// #define PN53X_REG_- 0x632D
	// #define PN53X_REG_- 0x632E
	PN53X_REG_CIU_RFlevelDet Register = 0x632F
	PN53X_REG_CIU_SIC_CLK_en Register = 0x6330
	PN53X_REG_CIU_Command    Register = 0x6331
	PN53X_REG_CIU_CommIEn    Register = 0x6332
	PN53X_REG_CIU_DivIEn     Register = 0x6333
	PN53X_REG_CIU_CommIrq    Register = 0x6334
	PN53X_REG_CIU_DivIrq     Register = 0x6335
	PN53X_REG_CIU_Error      Register = 0x6336
	PN53X_REG_CIU_Status1    Register = 0x6337
	PN53X_REG_CIU_Status2    Register = 0x6338
	PN53X_REG_CIU_FIFOData   Register = 0x6339
	PN53X_REG_CIU_FIFOLevel  Register = 0x633A
	PN53X_REG_CIU_WaterLevel Register = 0x633B
	PN53X_REG_CIU_Control    Register = 0x633C
	PN53X_REG_CIU_BitFraming Register = 0x633D
	PN53X_REG_CIU_Coll       Register = 0x633E

	PN53X_SFR_P3 Register = 0xFFB0

	PN53X_SFR_P3CFGA Register = 0xFFFC
	PN53X_SFR_P3CFGB Register = 0xFFFD
	PN53X_SFR_P7CFGA Register = 0xFFF4
	PN53X_SFR_P7CFGB Register = 0xFFF5
	PN53X_SFR_P7     Register = 0xFFF7
)

type registerInfo struct {
	val  uint16
	name string
	desc string
}

func mkRegisterInfo(val Register, name string, description string) registerInfo {
	return registerInfo{
		val:  uint16(val),
		name: name,
		desc: description,
	}
}

var pn53xRegisters map[Register]registerInfo = map[Register]registerInfo{
	PN53X_REG_CIU_Mode:        mkRegisterInfo(PN53X_REG_CIU_Mode, "PN53X_REG_CIU_Mode", "Defines general modes for transmitting and receiving"),
	PN53X_REG_CIU_TxMode:      mkRegisterInfo(PN53X_REG_CIU_TxMode, "PN53X_REG_CIU_TxMode", "Defines the transmission data rate and framing during transmission"),
	PN53X_REG_CIU_RxMode:      mkRegisterInfo(PN53X_REG_CIU_RxMode, "PN53X_REG_CIU_RxMode", "Defines the transmission data rate and framing during receiving"),
	PN53X_REG_CIU_TxControl:   mkRegisterInfo(PN53X_REG_CIU_TxControl, "PN53X_REG_CIU_TxControl", "Controls the logical behaviour of the antenna driver pins TX1 and TX2"),
	PN53X_REG_CIU_TxAuto:      mkRegisterInfo(PN53X_REG_CIU_TxAuto, "PN53X_REG_CIU_TxAuto", "Controls the settings of the antenna driver"),
	PN53X_REG_CIU_TxSel:       mkRegisterInfo(PN53X_REG_CIU_TxSel, "PN53X_REG_CIU_TxSel", "Selects the internal sources for the antenna driver"),
	PN53X_REG_CIU_RxSel:       mkRegisterInfo(PN53X_REG_CIU_RxSel, "PN53X_REG_CIU_RxSel", "Selects internal receiver settings"),
	PN53X_REG_CIU_RxThreshold: mkRegisterInfo(PN53X_REG_CIU_RxThreshold, "PN53X_REG_CIU_RxThreshold", "Selects thresholds for the bit decoder"),
	PN53X_REG_CIU_Demod:       mkRegisterInfo(PN53X_REG_CIU_Demod, "PN53X_REG_CIU_Demod", "Defines demodulator settings"),
	PN53X_REG_CIU_FelNFC1:     mkRegisterInfo(PN53X_REG_CIU_FelNFC1, "PN53X_REG_CIU_FelNFC1", "Defines the length of the valid range for the received frame"),
	PN53X_REG_CIU_FelNFC2:     mkRegisterInfo(PN53X_REG_CIU_FelNFC2, "PN53X_REG_CIU_FelNFC2", "Defines the length of the valid range for the received frame"),
	PN53X_REG_CIU_MifNFC:      mkRegisterInfo(PN53X_REG_CIU_MifNFC, "PN53X_REG_CIU_MifNFC", "Controls the communication in ISO/IEC 14443/MIFARE and NFC target mode at 106 kbit/s"),
	PN53X_REG_CIU_ManualRCV:   mkRegisterInfo(PN53X_REG_CIU_ManualRCV, "PN53X_REG_CIU_ManualRCV", "Allows manual fine tuning of the internal receiver"),
	PN53X_REG_CIU_TypeB:       mkRegisterInfo(PN53X_REG_CIU_TypeB, "PN53X_REG_CIU_TypeB", "Configure the ISO/IEC 14443 type B"),
	// PNREG (PN53X_REG_-, "Reserved"),
	// PNREG (PN53X_REG_-, "Reserved"),
	PN53X_REG_CIU_CRCResultMSB:   mkRegisterInfo(PN53X_REG_CIU_CRCResultMSB, "PN53X_REG_CIU_CRCResultMSB", "Shows the actual MSB values of the CRC calculation"),
	PN53X_REG_CIU_CRCResultLSB:   mkRegisterInfo(PN53X_REG_CIU_CRCResultLSB, "PN53X_REG_CIU_CRCResultLSB", "Shows the actual LSB values of the CRC calculation"),
	PN53X_REG_CIU_GsNOFF:         mkRegisterInfo(PN53X_REG_CIU_GsNOFF, "PN53X_REG_CIU_GsNOFF", "Selects the conductance of the antenna driver pins TX1 and TX2 for load modulation when own RF field is switched OFF"),
	PN53X_REG_CIU_ModWidth:       mkRegisterInfo(PN53X_REG_CIU_ModWidth, "PN53X_REG_CIU_ModWidth", "Controls the setting of the width of the Miller pause"),
	PN53X_REG_CIU_TxBitPhase:     mkRegisterInfo(PN53X_REG_CIU_TxBitPhase, "PN53X_REG_CIU_TxBitPhase", "Bit synchronization at 106 kbit/s"),
	PN53X_REG_CIU_RFCfg:          mkRegisterInfo(PN53X_REG_CIU_RFCfg, "PN53X_REG_CIU_RFCfg", "Configures the receiver gain and RF level"),
	PN53X_REG_CIU_GsNOn:          mkRegisterInfo(PN53X_REG_CIU_GsNOn, "PN53X_REG_CIU_GsNOn", "Selects the conductance of the antenna driver pins TX1 and TX2 for modulation, when own RF field is switched ON"),
	PN53X_REG_CIU_CWGsP:          mkRegisterInfo(PN53X_REG_CIU_CWGsP, "PN53X_REG_CIU_CWGsP", "Selects the conductance of the antenna driver pins TX1 and TX2 when not in modulation phase"),
	PN53X_REG_CIU_ModGsP:         mkRegisterInfo(PN53X_REG_CIU_ModGsP, "PN53X_REG_CIU_ModGsP", "Selects the conductance of the antenna driver pins TX1 and TX2 when in modulation phase"),
	PN53X_REG_CIU_TMode:          mkRegisterInfo(PN53X_REG_CIU_TMode, "PN53X_REG_CIU_TMode", "Defines settings for the internal timer"),
	PN53X_REG_CIU_TPrescaler:     mkRegisterInfo(PN53X_REG_CIU_TPrescaler, "PN53X_REG_CIU_TPrescaler", "Defines settings for the internal timer"),
	PN53X_REG_CIU_TReloadVal_hi:  mkRegisterInfo(PN53X_REG_CIU_TReloadVal_hi, "PN53X_REG_CIU_TReloadVal_hi", "Describes the 16-bit long timer reload value (Higher 8 bits)"),
	PN53X_REG_CIU_TReloadVal_lo:  mkRegisterInfo(PN53X_REG_CIU_TReloadVal_lo, "PN53X_REG_CIU_TReloadVal_lo", "Describes the 16-bit long timer reload value (Lower 8 bits)"),
	PN53X_REG_CIU_TCounterVal_hi: mkRegisterInfo(PN53X_REG_CIU_TCounterVal_hi, "PN53X_REG_CIU_TCounterVal_hi", "Describes the 16-bit long timer actual value (Higher 8 bits)"),
	PN53X_REG_CIU_TCounterVal_lo: mkRegisterInfo(PN53X_REG_CIU_TCounterVal_lo, "PN53X_REG_CIU_TCounterVal_lo", "Describes the 16-bit long timer actual value (Lower 8 bits)"),
	// PNREG (PN53X_REG_-, "Reserved"),
	PN53X_REG_CIU_TestSel1:     mkRegisterInfo(PN53X_REG_CIU_TestSel1, "PN53X_REG_CIU_TestSel1", "General test signals configuration"),
	PN53X_REG_CIU_TestSel2:     mkRegisterInfo(PN53X_REG_CIU_TestSel2, "PN53X_REG_CIU_TestSel2", "General test signals configuration and PRBS control"),
	PN53X_REG_CIU_TestPinEn:    mkRegisterInfo(PN53X_REG_CIU_TestPinEn, "PN53X_REG_CIU_TestPinEn", "Enables test signals output on pins."),
	PN53X_REG_CIU_TestPinValue: mkRegisterInfo(PN53X_REG_CIU_TestPinValue, "PN53X_REG_CIU_TestPinValue", "Defines the values for the 8-bit parallel bus when it is used as I/O bus"),
	PN53X_REG_CIU_TestBus:      mkRegisterInfo(PN53X_REG_CIU_TestBus, "PN53X_REG_CIU_TestBus", "Shows the status of the internal test bus"),
	PN53X_REG_CIU_AutoTest:     mkRegisterInfo(PN53X_REG_CIU_AutoTest, "PN53X_REG_CIU_AutoTest", "Controls the digital self-test"),
	PN53X_REG_CIU_Version:      mkRegisterInfo(PN53X_REG_CIU_Version, "PN53X_REG_CIU_Version", "Shows the CIU version"),
	PN53X_REG_CIU_AnalogTest:   mkRegisterInfo(PN53X_REG_CIU_AnalogTest, "PN53X_REG_CIU_AnalogTest", "Controls the pins AUX1 and AUX2"),
	PN53X_REG_CIU_TestDAC1:     mkRegisterInfo(PN53X_REG_CIU_TestDAC1, "PN53X_REG_CIU_TestDAC1", "Defines the test value for the TestDAC1"),
	PN53X_REG_CIU_TestDAC2:     mkRegisterInfo(PN53X_REG_CIU_TestDAC2, "PN53X_REG_CIU_TestDAC2", "Defines the test value for the TestDAC2"),
	PN53X_REG_CIU_TestADC:      mkRegisterInfo(PN53X_REG_CIU_TestADC, "PN53X_REG_CIU_TestADC", "Show the actual value of ADC I and Q"),
	// PNREG (PN53X_REG_-, "Reserved for tests"),
	// PNREG (PN53X_REG_-, "Reserved for tests"),
	// PNREG (PN53X_REG_-, "Reserved for tests"),
	PN53X_REG_CIU_RFlevelDet: mkRegisterInfo(PN53X_REG_CIU_RFlevelDet, "PN53X_REG_CIU_RFlevelDet", "Power down of the RF level detector"),
	PN53X_REG_CIU_SIC_CLK_en: mkRegisterInfo(PN53X_REG_CIU_SIC_CLK_en, "PN53X_REG_CIU_SIC_CLK_en", "Enables the use of secure IC clock on P34 / SIC_CLK"),
	PN53X_REG_CIU_Command:    mkRegisterInfo(PN53X_REG_CIU_Command, "PN53X_REG_CIU_Command", "Starts and stops the command execution"),
	PN53X_REG_CIU_CommIEn:    mkRegisterInfo(PN53X_REG_CIU_CommIEn, "PN53X_REG_CIU_CommIEn", "Control bits to enable and disable the passing of interrupt requests"),
	PN53X_REG_CIU_DivIEn:     mkRegisterInfo(PN53X_REG_CIU_DivIEn, "PN53X_REG_CIU_DivIEn", "Controls bits to enable and disable the passing of interrupt requests"),
	PN53X_REG_CIU_CommIrq:    mkRegisterInfo(PN53X_REG_CIU_CommIrq, "PN53X_REG_CIU_CommIrq", "Contains common CIU interrupt request flags"),
	PN53X_REG_CIU_DivIrq:     mkRegisterInfo(PN53X_REG_CIU_DivIrq, "PN53X_REG_CIU_DivIrq", "Contains miscellaneous interrupt request flags"),
	PN53X_REG_CIU_Error:      mkRegisterInfo(PN53X_REG_CIU_Error, "PN53X_REG_CIU_Error", "Error flags showing the error status of the last command executed"),
	PN53X_REG_CIU_Status1:    mkRegisterInfo(PN53X_REG_CIU_Status1, "PN53X_REG_CIU_Status1", "Contains status flags of the CRC, Interrupt Request System and FIFO buffer"),
	PN53X_REG_CIU_Status2:    mkRegisterInfo(PN53X_REG_CIU_Status2, "PN53X_REG_CIU_Status2", "Contain status flags of the receiver, transmitter and Data Mode Detector"),
	PN53X_REG_CIU_FIFOData:   mkRegisterInfo(PN53X_REG_CIU_FIFOData, "PN53X_REG_CIU_FIFOData", "In- and output of 64 byte FIFO buffer"),
	PN53X_REG_CIU_FIFOLevel:  mkRegisterInfo(PN53X_REG_CIU_FIFOLevel, "PN53X_REG_CIU_FIFOLevel", "Indicates the number of bytes stored in the FIFO"),
	PN53X_REG_CIU_WaterLevel: mkRegisterInfo(PN53X_REG_CIU_WaterLevel, "PN53X_REG_CIU_WaterLevel", "Defines the thresholds for FIFO under- and overflow warning"),
	PN53X_REG_CIU_Control:    mkRegisterInfo(PN53X_REG_CIU_Control, "PN53X_REG_CIU_Control", "Contains miscellaneous control bits"),
	PN53X_REG_CIU_BitFraming: mkRegisterInfo(PN53X_REG_CIU_BitFraming, "PN53X_REG_CIU_BitFraming", "Adjustments for bit oriented frames"),
	PN53X_REG_CIU_Coll:       mkRegisterInfo(PN53X_REG_CIU_Coll, "PN53X_REG_CIU_Coll", "Defines the first bit collision detected on the RF interface"),

	// SFR
	PN53X_SFR_P3CFGA: mkRegisterInfo(PN53X_SFR_P3CFGA, "PN53X_SFR_P3CFGA", "Port 3 configuration"),
	PN53X_SFR_P3CFGB: mkRegisterInfo(PN53X_SFR_P3CFGB, "PN53X_SFR_P3CFGB", "Port 3 configuration"),
	PN53X_SFR_P3:     mkRegisterInfo(PN53X_SFR_P3, "PN53X_SFR_P3", "Port 3 value"),
	PN53X_SFR_P7CFGA: mkRegisterInfo(PN53X_SFR_P7CFGA, "PN53X_SFR_P7CFGA", "Port 7 configuration"),
	PN53X_SFR_P7CFGB: mkRegisterInfo(PN53X_SFR_P7CFGB, "PN53X_SFR_P7CFGB", "Port 7 configuration"),
	PN53X_SFR_P7:     mkRegisterInfo(PN53X_SFR_P7, "PN53X_SFR_P7", "Port 7 value"),
}

func (pnd *Chip) regTrace(reg Register) {
	info := pn53xRegisters[reg]
	pnd.logger.Debugf("%v (%v)", info.name, info.desc)
}

func (pnd *Chip) writebackRegister() error {
	// TODO Check at each step (ReadRegister, WriteRegister) if we didn't exceed max supported frame length
	abtReadRegisterCmd := gonfc.BufferInit(PN53x_EXTENDED_FRAME__DATA_MAX_LEN)

	gonfc.BufferAppend(abtReadRegisterCmd, byte(ReadRegister))

	// First step, it looks for registers to be read before applying the requested mask
	pnd.wbTrigged = false
	for n := uint16(0); n < uint16(PN53X_CACHE_REGISTER_SIZE); n++ {
		if (pnd.wbMask[n] != 0x00) && (pnd.wbMask[n] != 0xff) {
			// This register needs to be read: mask is present but does not cover full data width (ie. mask != 0xff)
			var pn53xRegisterAddress uint16 = PN53X_CACHE_REGISTER_MIN_ADDRESS + n
			gonfc.BufferAppend(abtReadRegisterCmd, byte(pn53xRegisterAddress>>8))
			gonfc.BufferAppend(abtReadRegisterCmd, byte(pn53xRegisterAddress&0xff))
		}
	}

	if abtReadRegisterCmd.Len() > 1 {
		// It needs to read some registers
		abtRes := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
		// size_t szRes = sizeof(abtRes);
		// It transceives the previously constructed ReadRegister command
		// if ((res = pn53x_transceive(pnd, abtReadRegisterCmd, BUFFER_SIZE(abtReadRegisterCmd), abtRes, szRes, -1)) < 0) {
		if _, err := pnd.transceive(abtReadRegisterCmd.Bytes(), abtRes, -1); err != nil {
			return err
		}
		i := 0
		if pnd.chipType == PN533 {
			// PN533 prepends its answer by a status byte
			i = 1
		}
		for n := 0; n < PN53X_CACHE_REGISTER_SIZE; n++ {
			if (pnd.wbMask[n] != 0x00) && (pnd.wbMask[n] != 0xff) {
				pnd.wbData[n] = ((pnd.wbData[n] & pnd.wbMask[n]) | (abtRes[i] & (^pnd.wbMask[n])))
				if pnd.wbData[n] != abtRes[i] {
					// Requested value is different from read one
					pnd.wbMask[n] = 0xff // We can now apply whole data bits
				} else {
					pnd.wbMask[n] = 0x00 // We already have the right value
				}
				i++
			}
		}
	}
	// Now, the writeback-cache only has masks with 0xff, we can start to WriteRegister
	abtWriteRegisterCmd := gonfc.BufferInit(PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	gonfc.BufferAppend(abtWriteRegisterCmd, byte(WriteRegister))
	for n := uint16(0); n < uint16(PN53X_CACHE_REGISTER_SIZE); n++ {
		if pnd.wbMask[n] == 0xff {
			var pn53xRegisterAddress uint16 = uint16(PN53X_CACHE_REGISTER_MIN_ADDRESS + n)
			pnd.regTrace(Register(pn53xRegisterAddress))
			gonfc.BufferAppend(abtWriteRegisterCmd, byte(pn53xRegisterAddress>>8))
			gonfc.BufferAppend(abtWriteRegisterCmd, byte(pn53xRegisterAddress&0xff))
			gonfc.BufferAppend(abtWriteRegisterCmd, pnd.wbData[n])
			// This register is handled, we reset the mask to prevent
			pnd.wbMask[n] = 0x00
		}
	}

	if abtWriteRegisterCmd.Len() > 1 {
		// We need to write some registers
		if _, err := pnd.transceive(abtWriteRegisterCmd.Bytes(), nil, -1); err != nil {
			return err
		}
	}
	return nil
}

func (pnd *Chip) writeRegister(ui16RegisterAddress Register, ui8Value byte) error {
	abtCmd := []byte{
		byte(WriteRegister),
		byte(ui16RegisterAddress >> 8),
		byte(ui16RegisterAddress & 0xff),
		ui8Value,
	}

	pnd.regTrace(ui16RegisterAddress)
	_, err := pnd.transceive(abtCmd, nil, -1)
	return err
}

func (pnd *Chip) writeRegisterMask(ui16RegisterAddress Register, ui8SymbolMask byte, ui8Value byte) error {
	if (uint16(ui16RegisterAddress) < PN53X_CACHE_REGISTER_MIN_ADDRESS) || (uint16(ui16RegisterAddress) > PN53X_CACHE_REGISTER_MAX_ADDRESS) {
		if ui8SymbolMask == 0xff {
			return pnd.writeRegister(ui16RegisterAddress, ui8Value)
		}

		ui8CurrentValue, err := pnd.readRegister(ui16RegisterAddress)
		if err != nil {
			return err
		}
		ui8NewValue := ((ui8Value & ui8SymbolMask) | (ui8CurrentValue & (^ui8SymbolMask)))
		if ui8NewValue != ui8CurrentValue {
			return pnd.writeRegister(ui16RegisterAddress, ui8NewValue)
		}
	} else {
		// Write-back cache area
		internalAddress := uint16(ui16RegisterAddress) - PN53X_CACHE_REGISTER_MIN_ADDRESS
		pnd.wbData[internalAddress] = (pnd.wbData[internalAddress] & pnd.wbMask[internalAddress] & (^ui8SymbolMask)) | (ui8Value & ui8SymbolMask)
		pnd.wbMask[internalAddress] = pnd.wbMask[internalAddress] | ui8SymbolMask
		pnd.wbTrigged = true
	}
	return nil
}

func (pnd *Chip) readRegister(ui16RegisterAddress Register) (byte, error) {
	abtCmd := []byte{
		byte(ReadRegister),
		byte(ui16RegisterAddress >> 8),
		byte(ui16RegisterAddress & 0xff),
	}
	abtRegValue := make([]byte, 2)

	pnd.regTrace(ui16RegisterAddress)
	_, err := pnd.transceive(abtCmd, abtRegValue, -1)
	if err != nil {
		return 0, err
	}

	if pnd.chipType == PN533 {
		// PN533 prepends its answer by a status byte
		return abtRegValue[1], nil
	}
	return abtRegValue[0], nil
}

func (pnd *Chip) setTxBits(ui8Bits byte) error {
	// Test if we need to update the transmission bits register setting
	if pnd.ui8TxBits == ui8Bits {
		return nil
	}
	if err := pnd.writeRegisterMask(PN53X_REG_CIU_BitFraming, SYMBOL_TX_LAST_BITS, ui8Bits); err != nil {
		return err
	}
	// Store the new setting
	pnd.ui8TxBits = ui8Bits
	return nil
}
