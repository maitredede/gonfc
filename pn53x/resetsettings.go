package pn53x

import "github.com/maitredede/gonfc"

func (pnd *Chip) resetSettings() error {
	// pnd.logger.Debugf("resetSettings")
	pnd.ui8TxBits = 0
	// Reset the ending transmission bits register, it is unknown what the last tranmission used there
	if err := pnd.writeRegisterMask(PN53X_REG_CIU_BitFraming, SYMBOL_TX_LAST_BITS, 0x00); err != nil {
		return err
	}
	// Make sure we reset the CRC and parity to chip handling.
	if err := pnd.SetPropertyBool(gonfc.NP_HANDLE_CRC, true); err != nil {
		return err
	}
	if err := pnd.SetPropertyBool(gonfc.NP_HANDLE_PARITY, true); err != nil {
		return err
	}
	// Activate "easy framing" feature by default
	if err := pnd.SetPropertyBool(gonfc.NP_EASY_FRAMING, true); err != nil {
		return err
	}
	// Deactivate the CRYPTO1 cipher, it may could cause problems when still active
	if err := pnd.SetPropertyBool(gonfc.NP_ACTIVATE_CRYPTO1, false); err != nil {
		return err
	}
	return nil
}
