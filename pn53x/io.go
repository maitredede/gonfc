package pn53x

type IO interface {
	Send(data []byte, timeout int) (int, error)
	Receive(data []byte, timeout int) (int, error)
	MaxPacketSize() int
}
