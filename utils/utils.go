package utils

import (
	"encoding/binary"
)

func ConvertUnsigned4Bytes(n uint32) []byte {
	var b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, n)
	return b
}

func LenghtUnsigned4Bytes(b []byte) int {
	return int(binary.LittleEndian.Uint32(b[:]))
}
