package pn53x

import "time"

type IO interface {
	Send(data []byte, timeout time.Duration) (int, error)
	Receive(data []byte, timeout time.Duration) (int, error)
	MaxPacketSize() int
}
