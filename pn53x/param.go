package pn53x

type Param byte

const (
	PARAM_NONE         Param = 0x00
	PARAM_NAD_USED     Param = 0x01
	PARAM_DID_USED     Param = 0x02
	PARAM_AUTO_ATR_RES Param = 0x04
	PARAM_AUTO_RATS    Param = 0x10
	PARAM_14443_4_PICC Param = 0x20 /* Only for PN532 */
	PARAM_NFC_SECURE   Param = 0x20 /* Only for PN533 */
	PARAM_NO_AMBLE     Param = 0x40 /* Only for PN532 */
)
