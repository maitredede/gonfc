package gonfc

func ValidateModulation(pnd Device, mode Mode, nm Modulation) error {

	nmt, err := pnd.DeviceGetSupportedModulation(mode)
	if err != nil {
		return err
	}
	for _, i := range nmt {
		if i != nm.Type {
			continue
		}
		var nbr []BaudRate
		var err error
		if mode == N_INITIATOR {
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
		return NFC_EINVARG
	}
	return NFC_EINVARG
}
