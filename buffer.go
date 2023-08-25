package gonfc

import "bytes"

func BufferInit(size int) *bytes.Buffer {
	b := &bytes.Buffer{}
	b.Grow(size)
	return b
}

func BufferAppend(buffer *bytes.Buffer, value byte) {
	if e := buffer.WriteByte(value); e != nil {
		panic(e)
	}
}
