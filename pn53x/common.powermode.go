package pn53x

type PowerMode byte

const (
	PowerModeNormal PowerMode = iota
	PowerModePowerDown
	PowerModeLowVBat
)

func (c *chipCommon) PowerMode() PowerMode {
	return c.powerMode
}

func (c *chipCommon) SetPowerMode(mode PowerMode) {
	c.powerMode = mode
}
