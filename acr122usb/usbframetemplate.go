package acr122usb

var usbFrameTemplate []byte = []byte{
	PC_to_RDR_XfrBlock, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // CCID header
	0xff, 0x00, 0x00, 0x00, 0x00, // ADPU header
	0xd4, // PN532 direction
}
