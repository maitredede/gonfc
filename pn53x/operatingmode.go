package pn53x

type OperatingMode byte

const (
	OperatingModeIdle OperatingMode = iota
	OperatingModeInitiator
	OperatingModeTarget
)
