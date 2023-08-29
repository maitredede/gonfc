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

var strModulationTypeData map[ModulationType]string = map[ModulationType]string{
	NMT_ISO14443A:       "ISO/IEC 14443A",
	NMT_ISO14443B:       "ISO/IEC 14443-4B",
	NMT_ISO14443BI:      "ISO/IEC 14443-4B",
	NMT_ISO14443BICLASS: "ISO/IEC 14443-2B-3B iClass (Picopass)",
	NMT_ISO14443B2CT:    "ISO/IEC 14443-2B ASK CTx",
	NMT_ISO14443B2SR:    "ISO/IEC 14443-2B ST SRx",
	NMT_FELICA:          "FeliCa",
	NMT_JEWEL:           "Innovision Jewel",
	NMT_BARCODE:         "Thinfilm NFC Barcode",
	NMT_DEP:             "D.E.P.",
}

func StrModulationType(nmt ModulationType) string {
	if s, ok := strModulationTypeData[nmt]; ok {
		return s
	}
	return "???"
}
