package gonfc

// NFC ISO14443A tag (MIFARE) information. ISO14443aTarget mirrors
// nfc_iso14443a_info.
type ISO14443aTarget struct {
	Atqa   [2]byte
	Sak    byte
	UIDLen int // length of the Uid field
	UID    [10]byte
	AtsLen int // length of the ATS field
	// Maximal theoretical ATS is FSD-2, FSD=256 for FSDI=8 in RATS
	Ats  [254]byte // up to 254 bytes
	Baud BaudRate  // Baud rate
}

var _ Target = (*ISO14443aTarget)(nil)

// Type is always ISO14443A
func (t *ISO14443aTarget) Modulation() Modulation {
	return Modulation{ISO14443a, t.Baud}
}
