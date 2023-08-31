package pn53x

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maitredede/gonfc"
)

func (d *Chip) InitiatorPollTarget(modulations []gonfc.Modulation, pollnr byte, period byte) (*gonfc.NfcTarget, error) {
	if d.chipType == PN532 {
		return d.initiatorPollTargetPN532(modulations, pollnr, period)
	} else {
		return d.initiatorPollTargetOther(modulations, pollnr, period)
	}
}

func (d *Chip) initiatorPollTargetPN532(modulations []gonfc.Modulation, pollnr byte, period byte) (*gonfc.NfcTarget, error) {
	apttTargetTypes := make([]TargetType, 0)
	for _, nm := range modulations {
		ptt := pn53x_nm_to_ptt(nm)
		if ptt == PTT_UNDEFINED {
			d.lastError.Set(gonfc.NFC_EINVARG)
			return nil, d.lastError.Get()
		}
		if (d.bAutoIso14443_4) && (ptt == PTT_MIFARE) { // Hack to have ATS
			apttTargetTypes = append(apttTargetTypes, PTT_ISO14443_4A_106, PTT_MIFARE)
		} else {
			apttTargetTypes = append(apttTargetTypes, ptt)
		}
	}
	ntTargets, err := d.InAutoPoll(apttTargetTypes, pollnr, period, 0)
	if err != nil {
		return nil, err
	}
	if len(ntTargets) == 0 {
		return nil, nil
	}
	if len(ntTargets) > 2 {
		return nil, gonfc.NFC_ECHIP
	}
	last := ntTargets[len(ntTargets)-1]
	return last, nil
}

func (d *Chip) initiatorPollTargetOther(modulations []gonfc.Modulation, pollnr byte, period byte) (*gonfc.NfcTarget, error) {
	panic("TODO initiatorPollTarget for non-PN532 chips")
}

func (d *Chip) InAutoPoll(targetTypes []TargetType, pollnr byte, period byte, timeout time.Duration) ([]*gonfc.NfcTarget, error) {
	if d.chipType != PN532 {
		// This function is not supported by pn531 neither pn533
		d.lastError.Set(gonfc.NFC_EDEVNOTSUPP)
		return nil, d.lastError.Get()
	}

	// InAutoPoll frame looks like this { 0xd4, 0x60, 0x0f, 0x01, 0x00 } => { direction, command, pollnr, period, types... }
	szTxInAutoPoll := 3 + len(targetTypes)
	abtCmd := make([]byte, 3+15)
	abtCmd[0] = byte(InAutoPoll)
	abtCmd[1] = byte(pollnr)
	abtCmd[2] = byte(period)
	for i := 0; i < len(targetTypes); i++ {
		abtCmd[3+i] = byte(targetTypes[i])
	}

	abtRx := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	res, err := d.transceive(abtCmd[:szTxInAutoPoll], abtRx, timeout)
	if err != nil {
		return nil, err
	}

	if res == 0 {
		return nil, nil
	}
	szTargetFound := int(abtRx[0])
	if szTargetFound == 0 {
		return nil, nil
	}
	br := bytes.NewBuffer(abtRx[1:])
	pntTargets := make([]*gonfc.NfcTarget, 0, szTargetFound)
	for i := 0; i < szTargetFound; i++ {

		pttByte, err := br.ReadByte()
		if err != nil {
			panic(err)
		}
		ptt := TargetType(pttByte)
		nm := pn53x_ptt_to_nm(ptt)

		ln, err := br.ReadByte()
		if err != nil {
			panic(err)
		}

		data := make([]byte, ln)
		if dataR, err := br.Read(data); err != nil {
			panic(err)
		} else {
			if dataR != len(data) {
				panic(fmt.Errorf("read length mismatch read=%v expected=%v", dataR, len(data)))
			}
		}

		nti, err := decodeTargetData(data, d.chipType, nm.Type)
		if err != nil {
			return nil, err
		}

		target := &gonfc.NfcTarget{
			NM:  nm,
			NTI: nti,
		}
		pntTargets = append(pntTargets, target)
	}
	return pntTargets, nil
}
