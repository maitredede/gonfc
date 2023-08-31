package gonfc

func InitiatorPollTarget(device Device, modulations []Modulation, pollnr int, period int) (*NfcTarget, error) {
	return device.InitiatorPollTarget(modulations, byte(pollnr), byte(period))
}
