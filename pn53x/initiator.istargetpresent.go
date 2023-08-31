package pn53x

import (
	"errors"
	"fmt"

	"github.com/maitredede/gonfc"
)

// InitiatorTargetIsPresent
// chips/pn53x.c pn53x_initiator_target_is_present
func (pnd *Chip) InitiatorTargetIsPresent(nt *gonfc.NfcTarget) (bool, error) {
	// Check if there is a saved target
	if pnd.currentTarget == nil {
		pnd.logger.Debug("target_is_present(): no saved target")
		pnd.lastError.Set(gonfc.NFC_EINVARG)
		return false, pnd.lastError.Get()
	}

	// Check if the argument target nt is equals to current saved target
	if nt != nil && !pnd.currentTargetIs(nt) {
		pnd.logger.Debug("target_is_present(): another target")
		pnd.lastError.Set(gonfc.NFC_ETGRELEASED)
		return false, pnd.lastError.Get()
	}

	// Ping target
	retRES, retERR := false, gonfc.NFC_EDEVNOTSUPP
	switch pnd.currentTarget.NM.Type {
	case gonfc.NMT_ISO14443A:
		nai := pnd.currentTarget.NTI.NAI()
		if (nai.BtSak & 0x20) != 0 {
			retRES, retERR = pnd.isPresent_ISO14443A_4()
		} else if nai.AbtAtqa[0] == 0x00 && nai.AbtAtqa[1] == 0x44 && nai.BtSak == 0x00 {
			retRES, retERR = pnd.isPresent_ISO14443A_MFUL()
		} else if (nai.BtSak & 0x08) != 0 {
			retRES, retERR = pnd.isPresent_ISO14443A_MFC()
		} else {
			pnd.logger.Debug("target_is_present(): card type A not supported")
		}
		retERR = gonfc.NFC_EDEVNOTSUPP
	// case NMT_DEP:
	// 	ret = pn53x_DEP_is_present(pnd);
	// 	break;
	//   case NMT_FELICA:
	// 	ret = pn53x_Felica_is_present(pnd);
	// 	break;
	//   case NMT_JEWEL:
	// 	ret = pn53x_ISO14443A_Jewel_is_present(pnd);
	// 	break;
	//   case NMT_BARCODE:
	// 	ret = pn53x_ISO14443A_Barcode_is_present(pnd);
	// 	break;
	//   case NMT_ISO14443B:
	// 	ret = pn53x_ISO14443B_4_is_present(pnd);
	// 	break;
	//   case NMT_ISO14443BI:
	// 	ret = pn53x_ISO14443B_I_is_present(pnd);
	// 	break;
	//   case NMT_ISO14443B2SR:
	// 	ret = pn53x_ISO14443B_SR_is_present(pnd);
	// 	break;
	//   case NMT_ISO14443B2CT:
	// 	ret = pn53x_ISO14443B_CT_is_present(pnd);
	// 	break;
	//   case NMT_ISO14443BICLASS:
	// 	ret = pn53x_ISO14443B_ICLASS_is_present(pnd);
	// 	break;
	default:
		panic(fmt.Sprintf("TODO : implement target_is_present for modulationType %v", pnd.currentTarget.NM.Type))
	}
	if errors.Is(retERR, gonfc.NFC_ETGRELEASED) {
		pnd.currentTarget = nil
	}
	pnd.lastError.Set(retERR)
	return retRES, retERR
}

func (pnd *Chip) isPresent_ISO14443A_4() (bool, error) {
	panic("TODO isPresent_ISO14443A_4")
}

func (pnd *Chip) isPresent_ISO14443A_MFUL() (bool, error) {
	panic("TODO isPresent_ISO14443A_MFUL")
}

func (pnd *Chip) isPresent_ISO14443A_MFC() (bool, error) {
	panic("TODO isPresent_ISO14443A_MFC")
}
