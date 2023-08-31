package pn53x

import "github.com/maitredede/gonfc"

var (
	pn53xSupportedModulationAsTarget []gonfc.ModulationType = []gonfc.ModulationType{gonfc.NMT_ISO14443A, gonfc.NMT_FELICA, gonfc.NMT_DEP}
)

func (pnd *Chip) Init() error {
	// pnd.logger.Debug("Init")
	if err := pnd.decodeFirmwareVersion(); err != nil {
		return err
	}

	if pnd.supported_modulation_as_initiator == nil {
		pnd.supported_modulation_as_initiator = make([]gonfc.ModulationType, 0)
		if (pnd.btSupportByte & SUPPORT_ISO14443A) != 0 {
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443A)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_FELICA)
		}
		if (pnd.btSupportByte & SUPPORT_ISO14443B) != 0 {
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443B)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443BI)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443B2SR)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443B2CT)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_ISO14443BICLASS)
		}
		if pnd.chipType != PN531 {
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_JEWEL)
			pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_BARCODE)
		}
		pnd.supported_modulation_as_initiator = append(pnd.supported_modulation_as_initiator, gonfc.NMT_DEP)
	}
	if pnd.supported_modulation_as_target == nil {
		pnd.supported_modulation_as_target = pn53xSupportedModulationAsTarget
	}

	// CRC handling should be enabled by default as declared in nfc_device_new
	// which is the case by default for pn53x, so nothing to do here
	// Parity handling should be enabled by default as declared in nfc_device_new
	// which is the case by default for pn53x, so nothing to do here

	// We can't read these parameters, so we set a default config by using the SetParameters wrapper
	// Note: pn53x_SetParameters() will save the sent value in pnd->ui8Parameters cache
	if err := pnd.SetParameters(PARAM_AUTO_ATR_RES | PARAM_AUTO_RATS); err != nil {
		return err
	}

	if err := pnd.resetSettings(); err != nil {
		return err
	}
	return nil
}
