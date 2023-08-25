package gonfc

// NFC FeLiCa tag information
type FelicaTarget struct {
	Len     uint
	ResCode byte
	ID      [8]byte
	Pad     [8]byte
	SysCode [2]byte
	Baud    BaudRate
}

var _ Target = (*FelicaTarget)(nil)

// Type is always FELICA
func (t *FelicaTarget) Modulation() Modulation {
	return Modulation{Felica, t.Baud}
}
