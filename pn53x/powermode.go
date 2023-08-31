package pn53x

type PowerMode byte

const (
	PowerModeNormal PowerMode = iota
	PowerModePowerDown
	PowerModeLowVBat
)

func (c *Chip) PowerMode() PowerMode {
	return c.powerMode
}

func (c *Chip) SetPowerMode(mode PowerMode) {
	c.powerMode = mode
}
