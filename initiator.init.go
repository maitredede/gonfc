package gonfc

func InitiatorInit(pnd Device) error {
	// Drop the field for a while
	if err := pnd.SetPropertyBool(NP_ACTIVATE_FIELD, false); err != nil {
		return err
	}
	// Enable field so more power consuming cards can power themselves up
	if err := pnd.SetPropertyBool(NP_ACTIVATE_FIELD, true); err != nil {
		return err
	}
	// Let the device try forever to find a target/tag
	if err := pnd.SetPropertyBool(NP_INFINITE_SELECT, true); err != nil {
		return err
	}
	// Activate auto ISO14443-4 switching by default
	if err := pnd.SetPropertyBool(NP_AUTO_ISO14443_4, true); err != nil {
		return err
	}
	// Force 14443-A mode
	if err := pnd.SetPropertyBool(NP_FORCE_ISO14443_A, true); err != nil {
		return err
	}
	// Force speed at 106kbps
	if err := pnd.SetPropertyBool(NP_FORCE_SPEED_106, true); err != nil {
		return err
	}
	// Disallow invalid frame
	if err := pnd.SetPropertyBool(NP_ACCEPT_INVALID_FRAMES, false); err != nil {
		return err
	}
	// Disallow multiple frames
	if err := pnd.SetPropertyBool(NP_ACCEPT_MULTIPLE_FRAMES, false); err != nil {
		return err
	}

	return pnd.InitiatorInit()
}
