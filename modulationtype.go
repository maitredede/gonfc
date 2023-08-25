package gonfc

type ModulationType byte

const (
	NMT_ISO14443A ModulationType = iota + 1
	NMT_JEWEL
	NMT_ISO14443B
	NMT_ISO14443BI   // pre-ISO14443B aka ISO/IEC 14443 B' or Type B'
	NMT_ISO14443B2SR // ISO14443-2B ST SRx
	NMT_ISO14443B2CT // ISO14443-2B ASK CTx
	NMT_FELICA
	NMT_DEP
	NMT_BARCODE                                              // Thinfilm NFC Barcode
	NMT_ISO14443BICLASS                                      // HID iClass 14443B mode
	NMT_END_ENUM        ModulationType = NMT_ISO14443BICLASS // dummy for sizing - always should alias last
)

const (
	ISO14443a       = NMT_ISO14443A
	Jewel           = NMT_JEWEL
	ISO14443b       = NMT_ISO14443B
	ISO14443bi      = NMT_ISO14443BI
	ISO14443b2sr    = NMT_ISO14443B2SR
	ISO14443b2ct    = NMT_ISO14443B2CT
	Felica          = NMT_FELICA
	DEP             = NMT_DEP
	Barcode         = NMT_BARCODE
	ISO14443biClass = NMT_ISO14443BICLASS
)
