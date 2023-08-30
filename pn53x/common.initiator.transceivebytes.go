package pn53x

import (
	"time"

	"github.com/maitredede/gonfc"
)

func (pnd *chipCommon) InitiatorTransceiveBytes(pbtTx []byte, pbtRx []byte, timeout time.Duration) (int, error) {
	var szExtraTxLen int
	abtCmd := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	szTx := len(pbtTx)

	// We can not just send bytes without parity if while the PN53X expects we handled them
	if !pnd.bPar.Get() {
		pnd.lastError.Set(gonfc.NFC_EINVARG)
		return 0, pnd.lastError.Get()
	}

	// Copy the data into the command frame
	if pnd.bEasyFraming.Get() {
		abtCmd[0] = byte(InDataExchange)
		abtCmd[1] = 1 /* target number */
		// memcpy(abtCmd+2, pbtTx, szTx)
		for i := 0; i < szTx; i++ {
			abtCmd[i+2] = pbtTx[i]
		}
		szExtraTxLen = 2
	} else {
		abtCmd[0] = byte(InCommunicateThru)
		// memcpy(abtCmd+1, pbtTx, szTx)
		for i := 0; i < szTx; i++ {
			abtCmd[i+1] = pbtTx[i]
		}
		szExtraTxLen = 1
	}

	// To transfer command frames bytes we can not have any leading bits, reset this to zero
	if err := pnd.setTxBits(0); err != nil {
		pnd.lastError.Set(err)
		return 0, pnd.lastError.Get()
	}
	// Send the frame to the PN53X chip and get the answer
	// We have to give the amount of bytes + (the two command bytes 0xD4, 0x42)
	abtRx := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	var szRxLen int
	var err error
	if szRxLen, err = pnd.transceive(abtCmd[:szTx+szExtraTxLen], abtRx, timeout); err != nil {
		pnd.lastError.Set(err)
		return 0, pnd.lastError.Get()
	}

	//const size_t szRxLen = (size_t)res - 1;
	szRxLen = szRxLen - 1
	// if pbtRx != NULL {
	// 	if szRxLen > szRx {
	// 		pnd.logger.Errorf("Buffer size is too short: %v available(s), %v needed", szRx, szRxLen)
	// 		return NFC_EOVFLOW
	// 	}
	// 	// Copy the received bytes
	// 	memcpy(pbtRx, abtRx+1, szRxLen)
	// }
	for i := 0; i < szRxLen; i++ {
		pbtRx[i] = abtRx[i+1]
	}
	// Everything went successful, we return received bytes count
	return szRxLen, nil
}
