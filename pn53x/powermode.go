package pn53x

type PowerMode byte

const (
	PowerModeNormal PowerMode = iota
	PowerModePowerDown
	PowerModeLowVBat
)
