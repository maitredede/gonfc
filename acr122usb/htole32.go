package acr122usb

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"unsafe"
// )

// // https://stackoverflow.com/a/53286786

// var nativeEndian binary.ByteOrder

// func init() {
// 	buf := [2]byte{}
// 	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

// 	switch buf {
// 	case [2]byte{0xCD, 0xAB}:
// 		nativeEndian = binary.LittleEndian
// 	case [2]byte{0xAB, 0xCD}:
// 		nativeEndian = binary.BigEndian
// 	default:
// 		panic("Could not determine native endianness.")
// 	}
// }

// func htole32(val uint32) uint32 {
// 	if nativeEndian == binary.LittleEndian {
// 		return val
// 	}
// 	b := &bytes.Buffer{}
// 	if err := binary.Write(b, binary.LittleEndian, val); err != nil {
// 		panic(err)
// 	}
// 	bin := b.Bytes()
// 	var ret uint32
// 	ret = *(*uint32)(unsafe.Pointer(&bin[0]))
// 	return ret
// }
