package pn53x

import (
	"bytes"
	"fmt"

	"github.com/maitredede/gonfc"
)

func decodeTargetData(rawData []byte, chipType ChipType, nmt gonfc.ModulationType) (*gonfc.NfcTargetInfo, error) {

	switch nmt {
	case gonfc.NMT_ISO14443A:
		return decodeTargetDataISO14443A(rawData, chipType)
	// case NMT_ISO14443B:
	// case NMT_ISO14443BI:
	// case NMT_ISO14443B2SR:
	// case NMT_ISO14443BICLASS:
	// case NMT_ISO14443B2CT:
	// case NMT_FELICA:
	// case NMT_JEWEL:
	// case NMT_BARCODE:
	case gonfc.NMT_DEP:
		return nil, gonfc.NFC_ECHIP
	default:
		panic("TODO")
	}
}

func decodeTargetDataISO14443A(rawData []byte, chipType ChipType) (*gonfc.NfcTargetInfo, error) {
	pnti := &gonfc.NfcTargetInfo{}
	rawBuff := bytes.NewBuffer(rawData)
	nextByte := func() byte {
		if b, err := rawBuff.ReadByte(); err != nil {
			panic(err)
		} else {
			return b
		}
	}
	nextBlock := func(l int) []byte {
		b := make([]byte, l)
		n, err := rawBuff.Read(b)
		if err != nil {
			panic(err)
		}
		if n != l {
			panic(fmt.Errorf("requested=%v read=%v", l, n))
		}
		return b[:n]
	}

	// We skip the first byte: its the target number (Tg)
	//pbtRawData++;
	nextByte()

	// Somehow they switched the lower and upper ATQA bytes around for the PN531 chipset
	if chipType == PN531 {
		pnti.NAI().AbtAtqa[1] = nextByte()
		pnti.NAI().AbtAtqa[0] = nextByte()
	} else {
		pnti.NAI().AbtAtqa[0] = nextByte()
		pnti.NAI().AbtAtqa[1] = nextByte()
	}
	pnti.NAI().BtSak = nextByte()
	// Copy the NFCID1
	szUidLen := int(nextByte())
	pnti.NAI().SzUidLen = uint32(szUidLen)
	pbtUID := nextBlock(szUidLen)

	// Did we received an optional ATS (Smardcard ATR)
	if len(rawData) > szUidLen+5 {
		szAtsLen := int(nextByte()) - 1
		pnti.NAI().SzAtsLen = uint32(szAtsLen)
		abtAts := nextBlock(szAtsLen)
		for i := 0; i < szAtsLen; i++ {
			pnti.NAI().AbtAts[i] = abtAts[i]
		}
	} else {
		pnti.NAI().SzAtsLen = 0
	}

	// For PN531, strip CT (Cascade Tag) to retrieve and store the _real_ UID
	// (e.g. 0x8801020304050607 is in fact 0x01020304050607)
	if pnti.NAI().SzUidLen == 8 && pbtUID[0] == 0x88 {
		pnti.NAI().SzUidLen = 7
		for i := 0; i < 7; i++ {
			pnti.NAI().AbtUid[i] = pbtUID[i+1]
		}
		//      } else if ((pnti.NAI().szUidLen == 12) && (pbtUid[0] == 0x88) && (pbtUid[4] == 0x88)) {
	} else if pnti.NAI().SzUidLen > 10 {
		pnti.NAI().SzUidLen = 10
		for i := 0; i < 3; i++ {
			pnti.NAI().AbtUid[i] = pbtUID[i+1]
			pnti.NAI().AbtUid[i+3] = pbtUID[i+5]
		}
		for i := 0; i < 4; i++ {
			pnti.NAI().AbtUid[i+6] = pbtUID[i+8]
		}
	} else {
		// For PN532, PN533
		for i := 0; i < int(pnti.NAI().SzUidLen); i++ {
			pnti.NAI().AbtUid[i] = pbtUID[i]
		}
	}
	return pnti, nil
}
