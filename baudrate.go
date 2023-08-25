package gonfc

type BaudRate byte

const (
	NbrUndefined BaudRate = iota
	Nbr106
	Nbr212
	Nbr424
	Nbr847
)
