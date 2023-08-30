package pn53x

type Param byte

const (
	PARAM_NONE         Param = 0x00
	PARAM_NAD_USED     Param = 0x01
	PARAM_DID_USED     Param = 0x02
	PARAM_AUTO_ATR_RES Param = 0x04
	PARAM_AUTO_RATS    Param = 0x10
	PARAM_14443_4_PICC Param = 0x20 /* Only for PN532 */
	PARAM_NFC_SECURE   Param = 0x20 /* Only for PN533 */
	PARAM_NO_AMBLE     Param = 0x40 /* Only for PN532 */
)

func (pnd *chipCommon) SetParameters(parameter Param) error {
	// pnd.logger.Debugf("SetParameters")
	v := byte(parameter)
	abtCmd := []byte{byte(SetParameters), v}
	if _, err := pnd.transceive(abtCmd, nil, -1); err != nil {
		return err
	}
	pnd.ui8Parameters = v
	return nil
}

func (pnd *chipCommon) SetParametersEnable(ui8Parameter Param, bEnable bool) error {
	// pnd.logger.Debugf("SetParametersEnable")
	var ui8Value byte
	if bEnable {
		ui8Value = byte(pnd.ui8Parameters | byte(ui8Parameter))
	} else {
		ui8Value = byte(pnd.ui8Parameters & ^(byte(ui8Parameter)))
	}
	if ui8Value != pnd.ui8Parameters {
		return pnd.SetParameters(Param(ui8Value))
	}
	return nil
}
