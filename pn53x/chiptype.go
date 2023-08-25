package pn53x

type ChipType byte

const (
	PN53x  ChipType = 0x00 // Unknown PN53x chip type
	PN531  ChipType = 0x01
	PN532  ChipType = 0x02
	PN533  ChipType = 0x04
	RCS360 ChipType = 0x08
)
