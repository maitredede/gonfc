package acr122usb

import "github.com/maitredede/gonfc"

func (pnd *Acr122UsbDevice) InitiatorInit() error {
	return pnd.chip.InitiatorInit()
}

// TODO move to gonfc
func (pnd *Acr122UsbDevice) InitiatorListPassiveTargets(nm gonfc.Modulation) ([]gonfc.Target, error) {

	ant := make([]gonfc.Target, 0)

	pnd.lastError = nil
	// Let the reader only try once to find a tag
	bInfiniteSelect := pnd.InfiniteSelect
	if err := pnd.SetPropertyBool(gonfc.InfiniteSelect, false); err != nil {
		return ant, err
	}

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

// TODO move to gonfc
func prepateInitiatorData(nm gonfc.Modulation) []byte {
	switch nm.Type {
	case gonfc.NMT_ISO14443B:
		// Application Family Identifier (AFI) must equals 0x00 in order to wakeup all ISO14443-B PICCs (see ISO/IEC 14443-3)
		return []byte{0}
	case gonfc.NMT_ISO14443BI:
		// APGEN
		return []byte{0x01, 0x0b, 0x3f, 0x80}
	case gonfc.NMT_FELICA:
		// polling payload must be present (see ISO/IEC 18092 11.2.2.5)
		return []byte{0x00, 0xff, 0xff, 0x01, 0x00}
	case gonfc.NMT_ISO14443A:
		fallthrough
	case gonfc.NMT_ISO14443B2CT:
		fallthrough
	case gonfc.NMT_ISO14443B2SR:
		fallthrough
	case gonfc.NMT_ISO14443BICLASS:
		fallthrough
	case gonfc.NMT_JEWEL:
		fallthrough
	case gonfc.NMT_BARCODE:
		fallthrough
	case gonfc.NMT_DEP:
		return nil
	}
	panic("unknown modulation")
}

func (pnd *Acr122UsbDevice) InitiatorSelectPassiveTarget(nm gonfc.Modulation, initData []byte) (gonfc.Target, error) {
	szInitData := len(initData)

	var abtInit []byte
	maxAbt := max(12, szInitData)
	abtTmpInit := make([]byte, maxAbt)
	szInit := 0

	if err := pnd.DeviceValidateModulation(gonfc.N_INITIATOR, nm); err != nil {
		return nil, err
	}

	if szInitData == 0 {
		// Provide default values, if any
		initData = prepateInitiatorData(nm)
	} else if nm.Type == gonfc.NMT_ISO14443A {
		abtInit = abtTmpInit
		szInit = gonfc.ISO14443CascadeUID(initData, abtInit)
	} else {
		abtInit = abtTmpInit
		// memcpy(abtInit, pbtInitData, szInitData);
		for i := 0; i < szInitData; i++ {
			abtInit[i] = initData[i]
		}
		szInit = szInitData
	}

	//HAL(initiator_select_passive_target
	nt, err := pnd.chip.InitiatorSelectPassiveTarget(nm, abtInit[:szInit])

	return nt, err
}

func (pnd *Acr122UsbDevice) DeviceValidateModulation(mode gonfc.Mode, nm gonfc.Modulation) error {

	nmt, err := pnd.DeviceGetSupportedModulation(mode)
	if err != nil {
		return err
	}
	for _, i := range nmt {
		if i != nm.Type {
			continue
		}
		var nbr []gonfc.BaudRate
		var err error
		if mode == gonfc.N_INITIATOR {
			nbr, err = pnd.GetSupportedBaudRate(i)
			if err != nil {
				return err
			}
		} else {
			nbr, err = pnd.GetSupportedBaudRateTargetMode(i)
			if err != nil {
				return err
			}
		}
		for _, j := range nbr {
			if j == nm.BaudRate {
				return nil
			}
		}
		return gonfc.NFC_EINVARG
	}
	return gonfc.NFC_EINVARG
}

func (pnd *Acr122UsbDevice) InitiatorTransceiveBytes(tx, rx []byte, timeout int) (int, error) {
	return pnd.chip.InitiatorTransceiveBytes(tx, rx, timeout)
}

func (pnd *Acr122UsbDevice) InitiatorDeselectTarget() error {
	return pnd.chip.InitiatorDeselectTarget()
}

func (pnd *Acr122UsbDevice) DeviceGetSupportedModulation(mode gonfc.Mode) ([]gonfc.ModulationType, error) {
	return pnd.chip.GetSupportedModulation(mode)
}

func (pnd *Acr122UsbDevice) GetSupportedBaudRate(mt gonfc.ModulationType) ([]gonfc.BaudRate, error) {
	return pnd.chip.GetSupportedBaudRate(gonfc.N_INITIATOR, mt)
}

func (pnd *Acr122UsbDevice) GetSupportedBaudRateTargetMode(mt gonfc.ModulationType) ([]gonfc.BaudRate, error) {
	return pnd.chip.GetSupportedBaudRate(gonfc.N_TARGET, mt)
}
