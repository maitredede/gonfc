package acr122usb

// CCID Bulk-Out messages type
const (
	PC_to_RDR_IccPowerOn byte = 0x62
	PC_to_RDR_XfrBlock   byte = 0x6f

	RDR_to_PC_DataBlock byte = 0x80
)

// ISO 7816-4
const (
	SW1_More_Data_Available                     byte = 0x61
	SW1_Warning_with_NV_changed                 byte = 0x63
	PN53x_Specific_Application_Level_Error_Code byte = 0x7f
)

// APDUs instructions
const (
	APDU_GetAdditionnalData byte = 0xc0
)
