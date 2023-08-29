package acr122usb

import "github.com/maitredede/gonfc"

// TODO move to gonfc
func (pnd *Acr122UsbDevice) InitiatorListPassiveTargets(nm gonfc.Modulation) ([]gonfc.Target, error) {

	ant := make([]gonfc.Target, 0)

	pnd.LastError = nil
	// Let the reader only try once to find a tag
	bInfiniteSelect := pnd.InfiniteSelect
	if err := pnd.SetPropertyBool(gonfc.InfiniteSelect, false); err != nil {
		return ant, err
	}
	pnd.logger.Debugf("=== psv: NP_INFINITE_SELECT ok ===\n")

	pbtInitData := prepateInitiatorData(nm)

	var nt gonfc.Target
	var err error
	for {
		nt, err = pnd.InitiatorSelectPassiveTarget(nm, pbtInitData)
		if err != nil {
			pnd.logger.Warnf("TODO : handle error InitiatorSelectPassiveTarget: %v", err)
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
			pnd.logger.Warnf("TODO : handle error InitiatorDeselectTarget: %v", err)
		}
		// deselect has no effect on FeliCa, Jewel and Thinfilm cards so we'll stop after one...
		// ISO/IEC 14443 B' cards are polled at 100% probability so it's not possible to detect correctly two cards at the same time
		if (nm.Type == gonfc.NMT_FELICA) || (nm.Type == gonfc.NMT_JEWEL) || (nm.Type == gonfc.NMT_BARCODE) ||
			(nm.Type == gonfc.NMT_ISO14443BI) || (nm.Type == gonfc.NMT_ISO14443B2SR) || (nm.Type == gonfc.NMT_ISO14443B2CT) {
			break
		}
	}

	if bInfiniteSelect {
		if err := pnd.SetPropertyBool(gonfc.InfiniteSelect, true); err != nil {
			return ant, err
		}
	}

	return ant, nil
}
