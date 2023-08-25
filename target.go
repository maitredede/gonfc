package gonfc

type Target interface {
	Modulation() Modulation
}

type NfcTarget struct {
	NTI NfcTargetInfo
	NM  Modulation
}

//	typedef union {
//		nfc_iso14443a_info nai;
//		nfc_felica_info nfi;
//		nfc_iso14443b_info nbi;
//		nfc_iso14443bi_info nii;
//		nfc_iso14443b2sr_info nsi;
//		nfc_iso14443b2ct_info nci;
//		nfc_jewel_info nji;
//		nfc_dep_info ndi;
//		nfc_barcode_info nti; // "t" for Thinfilm, "b" already used
//		nfc_iso14443biclass_info nhi; // hid iclass / picopass - nii already used
//	  } nfc_target_info;
type NfcTargetInfo struct {
	val any
}

func (i *NfcTargetInfo) NAI() *NfcIso14443aInfo {
	return i.val.(*NfcIso14443aInfo)
}

func (i *NfcTargetInfo) NFI() *NfcFelicaInfo {
	return i.val.(*NfcFelicaInfo)
}

type NfcIso14443aInfo struct {
	abtAtqa  [2]byte
	btSak    byte
	szUidLen uint32
	abtUid   [10]byte
	szAtsLen uint32
	abtAts   [254]byte // Maximal theoretical ATS is FSD-2, FSD=256 for FSDI=8 in RATS
}

type NfcFelicaInfo struct {
	szLen      uint32
	btResCode  byte
	abtId      [8]byte
	abtPad     [8]byte
	abtSysCode [2]byte
}

var _ Target = (*NfcTarget)(nil)

func (t *NfcTarget) Modulation() Modulation {
	return t.NM
}
