package pn53x

import "github.com/maitredede/gonfc"

func (pnd *chipCommon) GetSupportedModulation(mode gonfc.Mode) ([]gonfc.ModulationType, error) {
	switch mode {
	case gonfc.N_TARGET:
		return pnd.supported_modulation_as_target, nil
	case gonfc.N_INITIATOR:
		return pnd.supported_modulation_as_initiator, nil
	}
	return nil, gonfc.NFC_EINVARG
}

type PNModulation byte

const (
	/** Undefined modulation */
	PM_UNDEFINED PNModulation = 0xff /* -1 */
	/** ISO14443-A (NXP MIFARE) http://en.wikipedia.org/wiki/MIFARE */
	PM_ISO14443A_106 PNModulation = 0x00
	/** JIS X 6319-4 (Sony Felica) http://en.wikipedia.org/wiki/FeliCa */
	PM_FELICA_212 PNModulation = 0x01
	/** JIS X 6319-4 (Sony Felica) http://en.wikipedia.org/wiki/FeliCa */
	PM_FELICA_424 PNModulation = 0x02
	/** ISO14443-B http://en.wikipedia.org/wiki/ISO/IEC_14443 (Not supported by PN531) */
	PM_ISO14443B_106 PNModulation = 0x03
	/** Jewel Topaz (Innovision Research & Development) (Not supported by PN531) */
	PM_JEWEL_106 PNModulation = 0x04
	/** Thinfilm NFC Barcode (Not supported by PN531) */
	PM_BARCODE_106 PNModulation = 0x05
	/** ISO14443-B http://en.wikipedia.org/wiki/ISO/IEC_14443 (Not supported by PN531 nor PN532) */
	PM_ISO14443B_212 PNModulation = 0x06
	/** ISO14443-B http://en.wikipedia.org/wiki/ISO/IEC_14443 (Not supported by PN531 nor PN532) */
	PM_ISO14443B_424 PNModulation = 0x07
	/** ISO14443-B http://en.wikipedia.org/wiki/ISO/IEC_14443 (Not supported by PN531 nor PN532) */
	PM_ISO14443B_847 PNModulation = 0x08
)

func pn53x_nm_to_pm(nm gonfc.Modulation) PNModulation {
	switch nm.Type {
	case gonfc.NMT_ISO14443A:
		return PM_ISO14443A_106

	case gonfc.NMT_ISO14443B:
	case gonfc.NMT_ISO14443BICLASS:
		switch nm.BaudRate {
		case gonfc.Nbr106:
			return PM_ISO14443B_106
		case gonfc.Nbr212:
			return PM_ISO14443B_212
		case gonfc.Nbr424:
			return PM_ISO14443B_424
		case gonfc.Nbr847:
			return PM_ISO14443B_847
		case gonfc.NbrUndefined:
			// Nothing to do...
			break
		}
		break

	case gonfc.NMT_JEWEL:
		return PM_JEWEL_106

	case gonfc.NMT_BARCODE:
		return PM_BARCODE_106

	case gonfc.NMT_FELICA:
		switch nm.BaudRate {
		case gonfc.Nbr212:
			return PM_FELICA_212
		case gonfc.Nbr424:
			return PM_FELICA_424
		case gonfc.Nbr106:
			fallthrough
		case gonfc.Nbr847:
			fallthrough
		case gonfc.NbrUndefined:
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
	case gonfc.NMT_DEP:
		// Nothing to do...
		break
	}
	return PM_UNDEFINED
}
