package pn53x

import (
	"time"
)

const ( // PN53X_REG_CIU_BitFraming
	SYMBOL_START_SEND   byte = 0x80
	SYMBOL_RX_ALIGN     byte = 0x70
	SYMBOL_TX_LAST_BITS byte = 0x07
)

const ( //   PN53X_REG_CIU_RxMode
	SYMBOL_RX_CRC_ENABLE byte = 0x80
	SYMBOL_RX_SPEED      byte = 0x70
	SYMBOL_RX_NO_ERROR   byte = 0x08
	SYMBOL_RX_MULTIPLE   byte = 0x04
	// RX_FRAMING follow same scheme than TX_FRAMING
	SYMBOL_RX_FRAMING byte = 0x03
)

const (
	// Registers and symbols masks used to covers parts within a register
	//
	//	PN53X_REG_CIU_TxMode
	SYMBOL_TX_CRC_ENABLE byte = 0x80
	SYMBOL_TX_SPEED      byte = 0x70
	// TX_FRAMING bits explanation:
	//   00 : ISO/IEC 14443A/MIFARE and Passive Communication mode 106 kbit/s
	//   01 : Active Communication mode
	//   10 : FeliCa and Passive Communication mode at 212 kbit/s and 424 kbit/s
	//   11 : ISO/IEC 14443B
	SYMBOL_TX_FRAMING byte = 0x03
)

const ( //   PN53X_REG_CIU_TxAuto
	SYMBOL_FORCE_100_ASK byte = 0x40
	SYMBOL_AUTO_WAKE_UP  byte = 0x20
	SYMBOL_INITIAL_RF_ON byte = 0x04
)

const (
	//   PN53X_REG_CIU_Status2
	SYMBOL_MF_CRYPTO1_ON byte = 0x08

	//   PN53X_REG_CIU_ManualRCV
	SYMBOL_PARITY_DISABLE byte = 0x10
)

const (
	// PN53X Support Byte flags
	SUPPORT_ISO14443A byte = 0x01
	SUPPORT_ISO14443B byte = 0x02
	SUPPORT_ISO18092  byte = 0x04
)

const (
	/**
	 * Start bytes, packet length, length checksum, direction, packet checksum and postamble are overhead
	 */
	// The TFI is considered part of the overhead
	PN53x_NORMAL_FRAME__DATA_MAX_LEN   int = 254
	PN53x_NORMAL_FRAME__OVERHEAD       int = 8
	PN53x_EXTENDED_FRAME__DATA_MAX_LEN int = 264
	PN53x_EXTENDED_FRAME__OVERHEAD     int = 11
	PN53x_ACK_FRAME__LEN               int = 6
)

func pn53x_int_to_timeout(ms int) byte {
	var res byte
	if ms > 0 {
		res = 0x10
		for i := 3280; i > 1; i /= 2 {
			if ms > i {
				break
			}
			res--
		}
	}
	return res
}

func pn53x_duration_to_timeout(t time.Duration) byte {
	ms := int(t.Milliseconds())
	return pn53x_int_to_timeout(ms)
}
