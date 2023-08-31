package gonfc

func InitiatorTargetIsPresent(device Device, target *NfcTarget) (bool, error) {
	return device.InitiatorTargetIsPresent(target)
}
