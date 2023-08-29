package gonfc

type BaudRate byte

const (
	NbrUndefined BaudRate = iota
	Nbr106
	Nbr212
	Nbr424
	Nbr847
)

func StrBaudRate(nbr BaudRate) string {
	switch nbr {
	case NbrUndefined:
		return "undefined baud rate"
	case Nbr106:
		return "106 kbps"
	case Nbr212:
		return "212 kbps"
	case Nbr424:
		return "424 kbps"
	case Nbr847:
		return "847 kbps"
	}
	return "???"
}
