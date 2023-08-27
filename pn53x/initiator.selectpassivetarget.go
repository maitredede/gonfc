package pn53x

import "github.com/maitredede/gonfc"

func (c *Chip) InitiatorSelectPassiveTarget(nm gonfc.Modulation, pbtInitData []byte) (gonfc.Target, error) {
	return c.InitiatorSelectPassiveTargetExt(nm, pbtInitData, 300)
}

func (c *Chip) InitiatorSelectPassiveTargetExt(nm gonfc.Modulation, pbtInitData []byte, timeout int) (gonfc.Target, error) {
	if nm.Type == gonfc.NMT_ISO14443BI || nm.Type == gonfc.NMT_ISO14443B2SR || nm.Type == gonfc.NMT_ISO14443B2CT || nm.Type == gonfc.NMT_ISO14443BICLASS {
		return c.initiatorSelectPassiveTargetExtIso(nm, pbtInitData, timeout)
	} else if nm.Type == gonfc.NMT_BARCODE {
		return c.initiatorSelectPassiveTargetExtBarcode(nm, pbtInitData, timeout)
	} else {
		return c.initiatorSelectPassiveTargetExtOther(nm, pbtInitData, timeout)
	}
}

func (pnd *Chip) initiatorSelectPassiveTargetExtIso(nm gonfc.Modulation, pbtInitData []byte, timeout int) (gonfc.Target, error) {
	abtTargetsData := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	// szTargetsData := len(abtTargetsData)
	var szTargetsData int

	if pnd.chipType == RCS360 {
		// TODO add support for RC-S360, at the moment it refuses to send raw frames without a first select
		pnd.lastError.Set(gonfc.NFC_ENOTIMPL)
		return nil, pnd.lastError.Get()
	}
	// No native support in InListPassiveTarget so we do discovery by hand
	if err := pnd.SetPropertyBool(gonfc.NP_FORCE_ISO14443_B, true); err != nil {
		return nil, err
	}
	if err := pnd.SetPropertyBool(gonfc.NP_FORCE_SPEED_106, true); err != nil {
		return nil, err
	}
	if err := pnd.SetPropertyBool(gonfc.NP_HANDLE_CRC, true); err != nil {
		return nil, err
	}
	if err := pnd.SetPropertyBool(gonfc.NP_EASY_FRAMING, false); err != nil {
		return nil, err
	}

	found := false
	for {
		if nm.Type == gonfc.NMT_ISO14443B2SR {
			// Some work to do before getting the UID...
			abtInitiate := []byte{0x06, 0x00}
			szInitiateLen := 2
			abtSelect := []byte{0x0e, 0x00}
			abtRx := make([]byte, 1)
			pbtInitData := []byte{0x0b}
			szInitData := 1
			if err := pnd.writeRegisterMask(PN53X_REG_CIU_TxAuto, 0xef, 0x07); err != nil { // Initial RFOn, Tx2 RFAutoEn, Tx1 RFAutoEn
				return nil, err
			}
			if err := pnd.writeRegisterMask(PN53X_REG_CIU_CWGsP, 0x3f, 0x3f); err != nil { // Conductance of the P-Driver
				return nil, err
			}
			if err := pnd.writeRegisterMask(PN53X_REG_CIU_ModGsP, 0x3f, 0x12); err != nil { // Driver P-output conductance for the time of modulation
				return nil, err
			}

			// Getting random Chip_ID
			if _, err := pnd.InitiatorTransceiveBytes(abtInitiate[:szInitiateLen], abtRx, timeout); err != nil {
				if err == gonfc.NFC_ERFTRANS && pnd.lastStatusByte == 0x01 { // Chip timeout
					continue
				}
				return nil, err
			}

			abtSelect[1] = abtRx[0]
			var err error
			if szTargetsData, err = pnd.InitiatorTransceiveBytes(abtSelect, abtRx, timeout); err != nil {
				return nil, err
			}

			if szTargetsData, err = pnd.InitiatorTransceiveBytes(pbtInitData[:szInitData], abtTargetsData, timeout); err != nil {
				if err == gonfc.NFC_ERFTRANS && pnd.lastStatusByte == 0x01 { // Chip timeout
					continue
				}
				return nil, err
			}
		} else if nm.Type == gonfc.NMT_ISO14443B2CT {
			panic("TODO")
			// Some work to do before getting the UID...
			// abtReqt := []byte{0x10}
			// pbtInitData = []byte{0x9f, 0xff, 0xff}
			// szInitData = 3

			//   // Getting product code / fab code & store it in output buffer after the serial nr we'll obtain later
			//   if ((res = pn53x_initiator_transceive_bytes(pnd, abtReqt, sizeof(abtReqt), abtTargetsData + 2, sizeof(abtTargetsData) - 2, timeout)) < 0) {
			// 	if ((res == NFC_ERFTRANS) && (CHIP_DATA(pnd)->last_status_byte == 0x01)) { // Chip timeout
			// 	  continue;
			// 	} else
			// 	  return res;
			//   }
			//   szTargetsData = (size_t)res;
			//   if ((res = pn53x_initiator_transceive_bytes(pnd, pbtInitData, szInitData, abtTargetsData, sizeof(abtTargetsData), timeout)) < 0) {
			// 	if ((res == NFC_ERFTRANS) && (CHIP_DATA(pnd)->last_status_byte == 0x01)) { // Chip timeout
			// 	  continue;
			// 	} else
			// 	  return res;
			//   }
			//   szTargetsData = (size_t)res;
			//   if (szTargetsData != 2)
			// 	return 0; // Target is not ISO14443B2CT
			//   uint8_t abtRead[] = { 0xC4 }; // Reading UID_MSB (Read address 4)
			//   if ((res = pn53x_initiator_transceive_bytes(pnd, abtRead, sizeof(abtRead), abtTargetsData + 4, sizeof(abtTargetsData) - 4, timeout)) < 0) {
			// 	return res;
			//   }
			//   szTargetsData = 6; // u16 UID_LSB, u8 prod code, u8 fab code, u16 UID_MSB
		} else if nm.Type == gonfc.NMT_ISO14443BICLASS {
			panic("TODO")
			//   pn53x_initiator_init_iclass_modulation(pnd);
			//   //
			//   // Some work to do before getting the UID...
			//   // send ICLASS_ACTIVATE_ALL command - will get timeout as we don't expect response
			//   uint8_t abtReqt[] = { 0x0a }; // iClass ACTIVATE_ALL
			//   uint8_t abtAnticol[11];
			//   if (pn53x_initiator_transceive_bytes(pnd, abtReqt, sizeof(abtReqt), NULL, 0, timeout) < 0) {
			// 	log_put(LOG_GROUP, LOG_CATEGORY, NFC_LOG_PRIORITY_DEBUG, "got expected timeout on iClass activate all");
			// 	//if ((res == NFC_ERFTRANS) && (CHIP_DATA(pnd)->last_status_byte == 0x01)) { // Chip timeout
			// 	//  continue;
			// 	//} else
			// 	//  return res;
			//   }
			//   // do select - returned anticol contains 'handle' for tag if present
			//   abtReqt[0] = 0x0c; // iClass SELECT
			//   abtAnticol[0] = 0x81; // iClass ANTICOL
			//   if ((res = pn53x_initiator_transceive_bytes(pnd, abtReqt, sizeof(abtReqt), &abtAnticol[1], sizeof(abtAnticol) - 1, timeout)) < 0) {
			// 	log_put(LOG_GROUP, LOG_CATEGORY, NFC_LOG_PRIORITY_DEBUG, "timeout on iClass anticol");
			// 	return res;
			//   }
			//   // write back anticol handle to get UID
			//   if ((res = pn53x_initiator_transceive_bytes(pnd, abtAnticol, 9, abtTargetsData, 10, timeout)) < 0) {
			// 	log_put(LOG_GROUP, LOG_CATEGORY, NFC_LOG_PRIORITY_DEBUG, "timeout on iClass get UID");
			// 	return res;
			//   }
			//   log_put(LOG_GROUP, LOG_CATEGORY, NFC_LOG_PRIORITY_DEBUG, "iClass raw UID: %02x %02x %02x %02x %02x %02x %02x %02x", abtTargetsData[0], abtTargetsData[1], abtTargetsData[2], abtTargetsData[3], abtTargetsData[4], abtTargetsData[5], abtTargetsData[6], abtTargetsData[7]);
			//   szTargetsData = 8;
			//   nttmp.nm = nm;
			//   if ((res = pn53x_decode_target_data(abtTargetsData, szTargetsData, CHIP_DATA(pnd)->type, nm.nmt, &(nttmp.nti))) < 0) {
			// 	return res;
			//   }
		} else {
			panic("TODO")

			//   if ((res = pn53x_initiator_transceive_bytes(pnd, pbtInitData, szInitData, abtTargetsData, sizeof(abtTargetsData), timeout)) < 0) {
			// 	if ((res == NFC_ERFTRANS) && (CHIP_DATA(pnd)->last_status_byte == 0x01)) { // Chip timeout
			// 	  continue;
			// 	} else
			// 	  return res;
			//   }
			//   szTargetsData = (size_t)res;
		}

		_ = szTargetsData
		_ = found
		panic("TODO")
		// nttmp.nm = nm;
		// if ((res = pn53x_decode_target_data(abtTargetsData, szTargetsData, CHIP_DATA(pnd)->type, nm.nmt, &(nttmp.nti))) < 0) {
		//   return res;
		// }
		// if (nm.Type == gonfc.NMT_ISO14443BI) {
		//   // Select tag
		//   uint8_t abtAttrib[6];
		//   memcpy(abtAttrib, abtTargetsData, sizeof(abtAttrib));
		//   abtAttrib[1] = 0x0f; // ATTRIB
		//   if ((res = pn53x_initiator_transceive_bytes(pnd, abtAttrib, sizeof(abtAttrib), NULL, 0, timeout)) < 0) {
		// 	return res;
		//   }
		//   szTargetsData = (size_t)res;
		// }
		// found = true;
		// if !pnd.bInfiniteSelect.Get() {
		// 	break
		// }
	}
	// if !found {
	// 	return nil, nil
	// }

	// _ = szTargetsData
	// panic("WIP")
}

func (c *Chip) initiatorSelectPassiveTargetExtBarcode(nm gonfc.Modulation, pbtInitData []byte, timeout int) (gonfc.Target, error) {
	panic("TODO")
}

func (pnd *Chip) initiatorSelectPassiveTargetExtOther(nm gonfc.Modulation, pbtInitData []byte, timeout int) (gonfc.Target, error) {

	abtTargetsData := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	//size_t  szTargetsData = sizeof(abtTargetsData);

	pm := pn53x_nm_to_pm(nm)
	if (PM_UNDEFINED == pm) || (gonfc.NbrUndefined == nm.BaudRate) {
		pnd.lastError.Set(gonfc.NFC_EINVARG)
		return nil, pnd.lastError.Get()
	}

	res, err := pnd.cmdInListPassiveTarget(pm, 1, pbtInitData, pbtInitData, abtTargetsData, timeout)
	if err != nil {
		return nil, err
	}
	_ = res
	panic("TODO")
	// if ((res = pn53x_InListPassiveTarget(pnd, pm, 1, pbtInitData, szInitData, abtTargetsData, &szTargetsData, timeout)) <= 0)
	//   return res;

	// if (szTargetsData <= 1) // For Coverity to know szTargetsData is always > 1 if res > 0
	//   return 0;

	// nttmp.nm = nm;
	// if ((res = pn53x_decode_target_data(abtTargetsData + 1, szTargetsData - 1, CHIP_DATA(pnd)->type, nm.nmt, &(nttmp.nti))) < 0) {
	//   return res;
	// }
	// if ((nm.nmt == NMT_ISO14443A) && (nm.nbr != NBR_106)) {
	//   uint8_t pncmd_inpsl[4] = { InPSL, 0x01 };
	//   pncmd_inpsl[2] = nm.nbr - 1;
	//   pncmd_inpsl[3] = nm.nbr - 1;
	//   if ((res = pn53x_transceive(pnd, pncmd_inpsl, sizeof(pncmd_inpsl), NULL, 0, 0)) < 0) {
	//     return res;
	//   }
	// }
}

func (pnd *Chip) cmdInListPassiveTarget(pmInitModulation PNModulation, szMaxTargets byte, pbtInitData []byte, pbtInitiatorData []byte, pbtTargetsData []byte, timeout int) (int, error) {
	szInitiatorData := len(pbtInitiatorData)
	abtCmd := make([]byte, 15)
	abtCmd[0] = byte(InListPassiveTarget)
	abtCmd[0] = szMaxTargets // MaxTg

	switch pmInitModulation {
	case PM_ISO14443A_106:
		fallthrough
	case PM_FELICA_212:
		fallthrough
	case PM_FELICA_424:
		// all gone fine.
		break
	case PM_ISO14443B_106:
		if (pnd.btSupportByte & SUPPORT_ISO14443B) == 0 {
			// Eg. Some PN532 doesn't support type B!
			pnd.lastError.Set(gonfc.NFC_EDEVNOTSUPP)
			return 0, pnd.lastError.Get()
		}
		break
	case PM_JEWEL_106:
		fallthrough
	case PM_BARCODE_106:
		if pnd.chipType == PN531 {
			// These modulations are not supported by pn531
			pnd.lastError.Set(gonfc.NFC_EDEVNOTSUPP)
			return 0, pnd.lastError.Get()
		}
		break
	case PM_ISO14443B_212:
		fallthrough
	case PM_ISO14443B_424:
		fallthrough
	case PM_ISO14443B_847:
		if (pnd.chipType != PN533) || (pnd.btSupportByte&SUPPORT_ISO14443B == 0) {
			// These modulations are not supported by pn531 neither pn532
			pnd.lastError.Set(gonfc.NFC_EDEVNOTSUPP)
			return 0, pnd.lastError.Get()
		}
		break
	case PM_UNDEFINED:
		pnd.lastError.Set(gonfc.NFC_EINVARG)
		return 0, pnd.lastError.Get()
	}
	abtCmd[2] = byte(pmInitModulation) // BrTy, the type of init modulation used for polling a passive tag

	// Set the optional initiator data (used for Felica, ISO14443B, Topaz Polling or for ISO14443A selecting a specific UID).
	if len(pbtInitiatorData) > 0 {
		// memcpy(abtCmd + 3, pbtInitiatorData, szInitiatorData);
		for i := 0; i < szInitiatorData; i++ {
			abtCmd[i+3] = pbtInitData[i]
		}
	}
	n, err := pnd.transceive(abtCmd[:3+szInitiatorData], pbtTargetsData, timeout)
	if err != nil {
		return 0, err
	}
	// int res = 0;
	//
	//	if ((res = pn53x_transceive(pnd, abtCmd, 3 + szInitiatorData, pbtTargetsData, *pszTargetsData, timeout)) < 0) {
	//	  return res;
	//	}
	//
	// *pszTargetsData = (size_t) res;
	// return pbtTargetsData[0];
	_ = n
	panic("TODO")
}