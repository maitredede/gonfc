package pn53x

func (pnd *Chip) setTxBits(ui8Bits byte) error {
	// Test if we need to update the transmission bits register setting
	if pnd.ui8TxBits == ui8Bits {
		return nil
	}
	if err := pnd.writeRegisterMask(PN53X_REG_CIU_BitFraming, SYMBOL_TX_LAST_BITS, ui8Bits); err != nil {
		return err
	}
	// Store the new setting
	pnd.ui8TxBits = ui8Bits
	return nil
}
