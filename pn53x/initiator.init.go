package pn53x

// InitiatorInit
// chips/pn53x.c pn53x_initiator_init
func (pnd *Chip) InitiatorInit() error {
	if err := pnd.resetSettings(); err != nil {
		return err
	}
	if pnd.samMode != SamModeNormal {
		if err := pnd.PN532SAMConfiguration(SamModeNormal, -1); err != nil {
			return err
		}
	}
	if err := pnd.writeRegisterMask(PN53X_REG_CIU_Control, SYMBOL_INITIATOR, 0x10); err != nil {
		return err
	}
	pnd.operatingMode = OperatingModeInitiator
	return nil
}
