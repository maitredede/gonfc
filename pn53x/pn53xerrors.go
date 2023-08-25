package pn53x

type pn53xError byte

const (
	/* PN53x specific errors */
	ETIMEOUT     pn53xError = 0x01
	ECRC         pn53xError = 0x02
	EPARITY      pn53xError = 0x03
	EBITCOUNT    pn53xError = 0x04
	EFRAMING     pn53xError = 0x05
	EBITCOLL     pn53xError = 0x06
	ESMALLBUF    pn53xError = 0x07
	EBUFOVF      pn53xError = 0x09
	ERFTIMEOUT   pn53xError = 0x0a
	ERFPROTO     pn53xError = 0x0b
	EOVHEAT      pn53xError = 0x0d
	EINBUFOVF    pn53xError = 0x0e
	EINVPARAM    pn53xError = 0x10
	EDEPUNKCMD   pn53xError = 0x12
	EINVRXFRAM   pn53xError = 0x13
	EMFAUTH      pn53xError = 0x14
	ENSECNOTSUPP pn53xError = 0x18 // PN533 only
	EBCC         pn53xError = 0x23
	EDEPINVSTATE pn53xError = 0x25
	EOPNOTALL    pn53xError = 0x26
	ECMD         pn53xError = 0x27
	ETGREL       pn53xError = 0x29
	ECID         pn53xError = 0x2a
	ECDISCARDED  pn53xError = 0x2b
	ENFCID3      pn53xError = 0x2c
	EOVCURRENT   pn53xError = 0x2d
	ENAD         pn53xError = 0x2e
)
