package pn53x

import "github.com/maitredede/gonfc"

var (
	supportedBaudRateFelica         = []gonfc.BaudRate{gonfc.Nbr424, gonfc.Nbr212}
	supportedBaudRateIso14443aPN532 = []gonfc.BaudRate{gonfc.Nbr424, gonfc.Nbr212, gonfc.Nbr106}
	supportedBaudRateIso14443aPN533 = []gonfc.BaudRate{gonfc.Nbr847, gonfc.Nbr424, gonfc.Nbr212, gonfc.Nbr106}
	supportedBaudRateIso14443bPN532 = []gonfc.BaudRate{gonfc.Nbr106}
	supportedBaudRateIso14443bPN533 = []gonfc.BaudRate{gonfc.Nbr847, gonfc.Nbr424, gonfc.Nbr212, gonfc.Nbr106}

	supportedBaudRateJewel   = []gonfc.BaudRate{gonfc.Nbr106}
	supportedBaudRateBarcode = []gonfc.BaudRate{gonfc.Nbr106}
	supportedBaudRateDep     = []gonfc.BaudRate{gonfc.Nbr424, gonfc.Nbr212, gonfc.Nbr106}
)

func (pnd *Chip) GetSupportedBaudRate(mode gonfc.Mode, nmt gonfc.ModulationType) ([]gonfc.BaudRate, error) {
	switch nmt {
	case gonfc.NMT_FELICA:
		return supportedBaudRateFelica, nil
	case gonfc.NMT_ISO14443A:
		{
			if pnd.chipType != PN533 || mode == gonfc.N_TARGET {
				return supportedBaudRateIso14443aPN532, nil
			}
			return supportedBaudRateIso14443aPN533, nil
		}
	case gonfc.NMT_ISO14443B:
		{
			if pnd.chipType != PN533 {
				return supportedBaudRateIso14443bPN532, nil
			} else {
				return supportedBaudRateIso14443bPN533, nil
			}
		}
	case gonfc.NMT_ISO14443BI:
		fallthrough
	case gonfc.NMT_ISO14443B2SR:
		fallthrough
	case gonfc.NMT_ISO14443B2CT:
		fallthrough
	case gonfc.NMT_ISO14443BICLASS:
		return supportedBaudRateIso14443bPN532, nil
	case gonfc.NMT_JEWEL:
		return supportedBaudRateJewel, nil
	case gonfc.NMT_BARCODE:
		return supportedBaudRateBarcode, nil
	case gonfc.NMT_DEP:
		return supportedBaudRateDep, nil
	}
	return nil, gonfc.NFC_EINVARG
}
