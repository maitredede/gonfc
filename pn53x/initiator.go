package pn53x

func (pnd *Chip) InitiatorDeselectTarget() error {
	return pnd.InDeselect(0) // 0 mean deselect all selected targets
}

func (pnd *Chip) InDeselect(target byte) error {
	if pnd.chipType == RCS360 {
		// We should do act here *only* if a target was previously selected
		abtStatus := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
		abtCmdGetStatus := []byte{byte(GetGeneralStatus)}
		szStatus, err := pnd.transceive(abtCmdGetStatus, abtStatus, -1)
		if err != nil {
			return err
		}
		if (szStatus < 3) || (abtStatus[2] == 0) {
			return nil
		}
		// No much choice what to deselect actually...
		abtCmdRcs360 := []byte{byte(InDeselect), 0x01, 0x01}
		_, err = pnd.transceive(abtCmdRcs360, nil, -1)
		return err
	}
	abtCmd := []byte{byte(InDeselect), target}
	_, err := pnd.transceive(abtCmd, nil, -1)
	return err
}
