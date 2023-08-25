package pn53x

type pn53xRFCI byte

const ( // Radio Field Configure Items   // Configuration Data length
	RFCI_FIELD                 pn53xRFCI = 0x01 //  1
	RFCI_TIMING                pn53xRFCI = 0x02 //  3
	RFCI_RETRY_DATA            pn53xRFCI = 0x04 //  1
	RFCI_RETRY_SELECT          pn53xRFCI = 0x05 //  3
	RFCI_ANALOG_TYPE_A_106     pn53xRFCI = 0x0A // 11
	RFCI_ANALOG_TYPE_A_212_424 pn53xRFCI = 0x0B //  8
	RFCI_ANALOG_TYPE_B         pn53xRFCI = 0x0C //  3
	RFCI_ANALOG_TYPE_14443_4   pn53xRFCI = 0x0D //  9
)
