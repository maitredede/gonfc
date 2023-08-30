package pn53x

import "github.com/maitredede/gonfc"

func BuildFrame(pbtData []byte) ([]byte, error) {
	pbtFrame := make([]byte, PN532_BUFFER_LEN)
	copy(pbtFrame, PN53X_PREAMBLE_AND_START)
	szData := len(pbtData)

	if szData <= PN53x_NORMAL_FRAME__DATA_MAX_LEN {
		// LEN - Packet length = data length (len) + checksum (1) + end of stream marker (1)
		pbtFrame[3] = byte(szData + 1)
		// LCS - Packet length checksum
		pbtFrame[4] = byte(256 - (szData + 1))
		// TFI
		pbtFrame[5] = 0xD4
		// DATA - Copy the PN53X command into the packet buffer
		//memcpy(pbtFrame + 6, pbtData, szData);
		for i := 0; i < szData; i++ {
			pbtFrame[i+6] = pbtData[i]
		}

		// DCS - Calculate data payload checksum
		btDCS := byte(256 - 0xD4)
		for szPos := 0; szPos < szData; szPos++ {
			btDCS -= pbtData[szPos]
		}
		pbtFrame[6+szData] = btDCS

		// 0x00 - End of stream marker
		pbtFrame[szData+7] = 0x00

		// (*pszFrame) = szData + PN53x_NORMAL_FRAME__OVERHEAD
		return pbtFrame[:szData+PN53x_NORMAL_FRAME__OVERHEAD], nil
	}
	if szData <= PN53x_EXTENDED_FRAME__DATA_MAX_LEN {
		// Extended frame marker
		pbtFrame[3] = 0xff
		pbtFrame[4] = 0xff
		// LENm
		pbtFrame[5] = byte((szData + 1) >> 8)
		// LENl
		pbtFrame[6] = byte((szData + 1) & 0xff)
		// LCS
		lcs := 256 - int((pbtFrame[5]+pbtFrame[6])&0xff)
		pbtFrame[7] = byte(lcs)
		// TFI
		pbtFrame[8] = 0xD4
		// DATA - Copy the PN53X command into the packet buffer
		// memcpy(pbtFrame+9, pbtData, szData)
		for i := 0; i < szData; i++ {
			pbtFrame[i+9] = pbtData[i]
		}

		// DCS - Calculate data payload checksum
		btDCS := byte(256 - 0xD4)
		for szPos := 0; szPos < szData; szPos++ {
			btDCS -= pbtData[szPos]
		}
		pbtFrame[9+szData] = btDCS

		// 0x00 - End of stream marker
		pbtFrame[szData+10] = 0x00

		return pbtFrame[:szData+PN53x_EXTENDED_FRAME__OVERHEAD], nil
	}
	// log_put(LOG_GROUP, LOG_CATEGORY, NFC_LOG_PRIORITY_ERROR, "We can't send more than %d bytes in a raw (requested: %" PRIdPTR ")", PN53x_EXTENDED_FRAME__DATA_MAX_LEN, szData);
	return nil, gonfc.NFC_ECHIP
}
