package gonfc

// InitiatorListPassiveTargets List passive or emulated tags
// nfc.c nfc_initiator_list_passive_targets
func InitiatorListPassiveTargets(pnd Device, nm Modulation) ([]*NfcTarget, error) {

	ant := make([]*NfcTarget, 0)

	pnd.SetLastError(nil)
	// Let the reader only try once to find a tag
	bInfiniteSelect := pnd.GetInfiniteSelect()
	if err := pnd.SetPropertyBool(InfiniteSelect, false); err != nil {
		return ant, err
	}
	pnd.Logger().Debugf("=== psv: NP_INFINITE_SELECT ok ===\n")

	pbtInitData := PrepateInitiatorData(nm)

	var nt *NfcTarget
	var err error
	for {
		nt, err = InitiatorSelectPassiveTarget(pnd, nm, pbtInitData)
		if err != nil {
			pnd.Logger().Warnf("TODO : handle error InitiatorSelectPassiveTarget: %v", err)
			break
		}

		var seen bool
		for _, t := range ant {
			if nt == t {
				seen = true
				break
			}
		}
		if seen {
			break
		}

		ant = append(ant, nt)

		if err := pnd.InitiatorDeselectTarget(); err != nil {
			pnd.Logger().Warnf("TODO : handle error InitiatorDeselectTarget: %v", err)
		}
		// deselect has no effect on FeliCa, Jewel and Thinfilm cards so we'll stop after one...
		// ISO/IEC 14443 B' cards are polled at 100% probability so it's not possible to detect correctly two cards at the same time
		if (nm.Type == NMT_FELICA) || (nm.Type == NMT_JEWEL) || (nm.Type == NMT_BARCODE) ||
			(nm.Type == NMT_ISO14443BI) || (nm.Type == NMT_ISO14443B2SR) || (nm.Type == NMT_ISO14443B2CT) {
			break
		}
	}

	if bInfiniteSelect {
		if err := pnd.SetPropertyBool(InfiniteSelect, true); err != nil {
			return ant, err
		}
	}

	return ant, nil
}

func PrepateInitiatorData(nm Modulation) []byte {
	switch nm.Type {
	case NMT_ISO14443B:
		// Application Family Identifier (AFI) must equals 0x00 in order to wakeup all ISO14443-B PICCs (see ISO/IEC 14443-3)
		return []byte{0}
	case NMT_ISO14443BI:
		// APGEN
		return []byte{0x01, 0x0b, 0x3f, 0x80}
	case NMT_FELICA:
		// polling payload must be present (see ISO/IEC 18092 11.2.2.5)
		return []byte{0x00, 0xff, 0xff, 0x01, 0x00}
	case NMT_ISO14443A:
		fallthrough
	case NMT_ISO14443B2CT:
		fallthrough
	case NMT_ISO14443B2SR:
		fallthrough
	case NMT_ISO14443BICLASS:
		fallthrough
	case NMT_JEWEL:
		fallthrough
	case NMT_BARCODE:
		fallthrough
	case NMT_DEP:
		return nil
	}
	panic("unknown modulation")
}
