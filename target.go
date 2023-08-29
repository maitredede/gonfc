package gonfc

import (
	"fmt"
	"strings"

	"github.com/maitredede/gonfc/utils"
)

type NfcTarget struct {
	NTI *NfcTargetInfo
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
	if i.val == nil {
		i.val = &NfcIso14443aInfo{}
	}
	return i.val.(*NfcIso14443aInfo)
}

func (i *NfcTargetInfo) NFI() *NfcFelicaInfo {
	if i.val == nil {
		i.val = &NfcFelicaInfo{}
	}
	return i.val.(*NfcFelicaInfo)
}

func (i *NfcTargetInfo) NDI() *NfcDepInfo {
	if i.val == nil {
		i.val = &NfcDepInfo{}
	}
	return i.val.(*NfcDepInfo)
}

type NfcIso14443aInfo struct {
	AbtAtqa  [2]byte
	BtSak    byte
	SzUidLen uint32
	AbtUid   [10]byte
	SzAtsLen uint32
	AbtAts   [254]byte // Maximal theoretical ATS is FSD-2, FSD=256 for FSDI=8 in RATS
}

func (i *NfcIso14443aInfo) UID() []byte {
	return i.AbtUid[:i.SzUidLen]
}

func (i *NfcIso14443aInfo) ATS() []byte {
	return i.AbtAts[:i.SzAtsLen]
}

type NfcFelicaInfo struct {
	szLen      uint32
	btResCode  byte
	abtId      [8]byte
	abtPad     [8]byte
	abtSysCode [2]byte
}

type NfcDepInfo struct {
	/** NFCID3 */
	AbtNFCID3 [10]byte
	/** DID */
	BtDID byte
	/** Supported send-bit rate */
	BtBS byte
	/** Supported receive-bit rate */
	BtBR byte
	/** Timeout value */
	BtTO byte
	/** PP Parameters */
	BtPP byte
	/** General Bytes */
	AbtGB [48]byte
	SzGB  int
	/** DEP mode */
	Ndm DepMode
}

func (i *NfcDepInfo) GB() []byte {
	return i.AbtGB[:i.SzGB]
}

func (t *NfcTarget) String() string {
	sb := strings.Builder{}
	sb.WriteString(StrModulationType(t.NM.Type))
	modeStr := ""
	if t.NM.Type == NMT_DEP {
		if t.NTI.NDI().Ndm == NDM_ACTIVE {
			modeStr = "active mode"
		} else {
			modeStr = "passive mode"
		}
	}
	sb.WriteString(fmt.Sprintf(" (%s%s)", StrBaudRate(t.NM.BaudRate), modeStr))
	sb.WriteString(" target:\n")

	switch t.NM.Type {
	case NMT_ISO14443A:
		sb.WriteString(t.StringISO14443A(false))
	default:
		return fmt.Sprintf("TODO %s", sb.String())
	}
	return sb.String()
}

func (t *NfcTarget) StringISO14443A(verbose bool) string {
	nai := t.NTI.NAI()
	sb := strings.Builder{}
	sb.WriteString("    ATQA (SENS_RES): ")
	sb.WriteString(utils.ToHexString(nai.AbtAtqa[:]) + "\n")
	if verbose {
		sb.WriteString("* UID size: ")
		switch nai.AbtAtqa[1] & 0xc0 >> 6 {
		case 0:
			sb.WriteString("single\n")
		case 1:
			sb.WriteString("double\n")
		case 2:
			sb.WriteString("triple\n")
		case 3:
			sb.WriteString("RFU\n")
		}
		sb.WriteString("* bit frame anticollision ")
		switch nai.AbtAtqa[1] & 0x1f {
		case 0x01:
			fallthrough
		case 0x02:
			fallthrough
		case 0x04:
			fallthrough
		case 0x08:
			fallthrough
		case 0x10:
			sb.WriteString("supported\n")
			break
		default:
			sb.WriteString("not supported\n")
			break
		}
	}
	var v string
	if nai.AbtUid[0] == 0x08 {
		v = "3"
	} else {
		v = "1"
	}
	sb.WriteString(fmt.Sprintf("       UID (NFCID%s): ", v))
	sb.WriteString(utils.ToHexString(nai.UID()) + "\n")
	if verbose {
		if nai.AbtUid[0] == 0x08 {
			sb.WriteString("* Random UID\n")
		}
	}
	sb.WriteString("      SAK (SEL_RES): ")
	sb.WriteString(utils.ToHexString([]byte{nai.BtSak}) + "\n")
	if verbose {
		if (nai.BtSak & SAK_UID_NOT_COMPLETE) != 0 {
			sb.WriteString("* Warning! Cascade bit set: UID not complete\n")
		}
		if (nai.BtSak & SAK_ISO14443_4_COMPLIANT) != 0 {
			sb.WriteString("* Compliant with ISO/IEC 14443-4\n")
		} else {
			sb.WriteString("* Not compliant with ISO/IEC 14443-4\n")
		}
		if (nai.BtSak & SAK_ISO18092_COMPLIANT) != 0 {
			sb.WriteString("* Compliant with ISO/IEC 18092\n")
		} else {
			sb.WriteString("* Not compliant with ISO/IEC 18092\n")
		}
	}

	if nai.SzAtsLen > 0 {
		sb.WriteString("                ATS: ")
		sb.WriteString(utils.ToHexString(nai.ATS()) + "\n")
	}

	if nai.SzAtsLen > 0 && verbose {
		sb.WriteString("TODO : ATS decode")
	}
	if verbose {
		sb.WriteString("TODO : Fingerprinting based on MIFARE type Identification Procedure")
	}
	return sb.String()
}
