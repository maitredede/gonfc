package gonfc

// DepMode NFC D.E.P. (Data Exchange Protocol) active/passive mode
type DepMode byte

const (
	NDM_UNDEFINED DepMode = iota
	NDM_PASSIVE
	NDM_ACTIVE
)
