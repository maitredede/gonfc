package pn53x

import "github.com/maitredede/gonfc"

// TargetType NFC target type
// chips/pn53x.h pn53x_target_type
type TargetType byte

const (
	/** Undefined target type */
	PTT_UNDEFINED TargetType = 0xff /*-1*/
	/** Generic passive 106 kbps (ISO/IEC14443-4A, mifare, DEP) */
	PTT_GENERIC_PASSIVE_106 TargetType = 0x00
	/** Generic passive 212 kbps (FeliCa, DEP) */
	PTT_GENERIC_PASSIVE_212 TargetType = 0x01
	/** Generic passive 424 kbps (FeliCa, DEP) */
	PTT_GENERIC_PASSIVE_424 TargetType = 0x02
	/** Passive 106 kbps ISO/IEC14443-4B */
	PTT_ISO14443_4B_106 TargetType = 0x03
	/** Innovision Jewel tag */
	PTT_JEWEL_106 TargetType = 0x04
	/** Mifare card */
	PTT_MIFARE TargetType = 0x10
	/** FeliCa 212 kbps card */
	PTT_FELICA_212 TargetType = 0x11
	/** FeliCa 424 kbps card */
	PTT_FELICA_424 TargetType = 0x12
	/** Passive 106 kbps ISO/IEC 14443-4A */
	PTT_ISO14443_4A_106 TargetType = 0x20
	/** Passive 106 kbps ISO/IEC 14443-4B with TCL flag */
	PTT_ISO14443_4B_TCL_106 TargetType = 0x23
	/** DEP passive 106 kbps */
	PTT_DEP_PASSIVE_106 TargetType = 0x40
	/** DEP passive 212 kbps */
	PTT_DEP_PASSIVE_212 TargetType = 0x41
	/** DEP passive 424 kbps */
	PTT_DEP_PASSIVE_424 TargetType = 0x42
	/** DEP active 106 kbps */
	PTT_DEP_ACTIVE_106 TargetType = 0x80
	/** DEP active 212 kbps */
	PTT_DEP_ACTIVE_212 TargetType = 0x81
	/** DEP active 424 kbps */
	PTT_DEP_ACTIVE_424 TargetType = 0x82
)

func pn53x_nm_to_ptt(nm gonfc.Modulation) TargetType {
	switch nm.Type {
	case gonfc.NMT_ISO14443A:
		return PTT_MIFARE
	// return PTT_ISO14443_4A_106;

	case gonfc.NMT_ISO14443B:
		fallthrough
	case gonfc.NMT_ISO14443BICLASS:
		switch nm.BaudRate {
		case gonfc.Nbr106:
			return PTT_ISO14443_4B_106

		case gonfc.NbrUndefined:
			fallthrough
		case gonfc.Nbr212:
			fallthrough
		case gonfc.Nbr424:
			fallthrough
		case gonfc.Nbr847:
			// Nothing to do...
			break
		}
		break

	case gonfc.NMT_JEWEL:
		return PTT_JEWEL_106

	case gonfc.NMT_FELICA:
		switch nm.BaudRate {
		case gonfc.Nbr212:
			return PTT_FELICA_212

		case gonfc.Nbr424:
			return PTT_FELICA_424

		case gonfc.NbrUndefined:
			fallthrough
		case gonfc.Nbr106:
			fallthrough
		case gonfc.Nbr847:
			// Nothing to do...
			break
		}
		break

	case gonfc.NMT_ISO14443BI:
		fallthrough
	case gonfc.NMT_ISO14443B2SR:
		fallthrough
	case gonfc.NMT_ISO14443B2CT:
		fallthrough
	case gonfc.NMT_BARCODE:
		fallthrough
	case gonfc.NMT_DEP:
		// Nothing to do...
		break
	}
	return PTT_UNDEFINED
}

func pn53x_ptt_to_nm(ptt TargetType) gonfc.Modulation {
	switch ptt {
	case PTT_GENERIC_PASSIVE_106:
		fallthrough
	case PTT_GENERIC_PASSIVE_212:
		fallthrough
	case PTT_GENERIC_PASSIVE_424:
		fallthrough
	case PTT_UNDEFINED:
		// XXX This should not happend, how handle it cleanly ?
		break

	case PTT_MIFARE:
		fallthrough
	case PTT_ISO14443_4A_106:
		return gonfc.Modulation{Type: gonfc.NMT_ISO14443A, BaudRate: gonfc.Nbr106}

	case PTT_ISO14443_4B_106:
		fallthrough
	case PTT_ISO14443_4B_TCL_106:
		return gonfc.Modulation{Type: gonfc.NMT_ISO14443B, BaudRate: gonfc.Nbr106}

	case PTT_JEWEL_106:
		return gonfc.Modulation{Type: gonfc.NMT_JEWEL, BaudRate: gonfc.Nbr106}

	case PTT_FELICA_212:
		return gonfc.Modulation{Type: gonfc.NMT_FELICA, BaudRate: gonfc.Nbr212}

	case PTT_FELICA_424:
		return gonfc.Modulation{Type: gonfc.NMT_FELICA, BaudRate: gonfc.Nbr424}

	case PTT_DEP_PASSIVE_106:
		fallthrough
	case PTT_DEP_ACTIVE_106:
		return gonfc.Modulation{Type: gonfc.NMT_DEP, BaudRate: gonfc.Nbr106}

	case PTT_DEP_PASSIVE_212:
		fallthrough
	case PTT_DEP_ACTIVE_212:
		return gonfc.Modulation{Type: gonfc.NMT_DEP, BaudRate: gonfc.Nbr212}

	case PTT_DEP_PASSIVE_424:
		fallthrough
	case PTT_DEP_ACTIVE_424:
		return gonfc.Modulation{Type: gonfc.NMT_DEP, BaudRate: gonfc.Nbr424}
	}
	// We should never be here, this line silent compilation warning
	return gonfc.Modulation{Type: gonfc.NMT_ISO14443A, BaudRate: gonfc.Nbr106}
}
