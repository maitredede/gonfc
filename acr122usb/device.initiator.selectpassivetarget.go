package acr122usb

import "github.com/maitredede/gonfc"

func (pnd *Acr122UsbDevice) InitiatorSelectPassiveTarget(nm gonfc.Modulation, initData []byte) (*gonfc.NfcTarget, error) {
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

	pnd.logger.Debugf("=== psv: select_before_HAL ===")
	//HAL(initiator_select_passive_target
	nt, err := pnd.chip.InitiatorSelectPassiveTarget(nm, abtInit[:szInit])

	return nt, err
}
