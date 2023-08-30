package gonfc

// InitiatorSelectPassiveTarget Select a passive or emulated tag
// nfc.c nfc_initiator_list_passive_targets
func InitiatorSelectPassiveTarget(pnd Device, nm Modulation, initData []byte) (*NfcTarget, error) {
	szInitData := len(initData)

	var abtInit []byte
	maxAbt := max(12, szInitData)
	abtTmpInit := make([]byte, maxAbt)
	szInit := 0

	if err := ValidateModulation(pnd, N_INITIATOR, nm); err != nil {
		return nil, err
	}

	if szInitData == 0 {
		// Provide default values, if any
		initData = PrepateInitiatorData(nm)
	} else if nm.Type == NMT_ISO14443A {
		abtInit = abtTmpInit
		szInit = ISO14443CascadeUID(initData, abtInit)
	} else {
		abtInit = abtTmpInit
		// memcpy(abtInit, pbtInitData, szInitData);
		for i := 0; i < szInitData; i++ {
			abtInit[i] = initData[i]
		}
		szInit = szInitData
	}

	//HAL(initiator_select_passive_target
	nt, err := pnd.InitiatorSelectPassiveTarget(nm, abtInit[:szInit])

	return nt, err
}
