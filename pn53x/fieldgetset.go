package pn53x

type BoolFieldGetSet interface {
	Get() bool
	Set(v bool)
}

type ByteFieldGetSet interface {
	Get() byte
	Set(v byte)
}

type ErrorFieldGetSet interface {
	Get() error
	Set(v error)
}
