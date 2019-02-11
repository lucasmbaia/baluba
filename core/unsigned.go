package core

import (
	"encoding/binary"
)

type BalulaFile struct {
	Version	  uint32
	Hostname  string
	Option	  string
	FileName  string
	File	  []byte
}

func ConvertUnsigned4Bytes(n uint32) []byte {
	var b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, n)
	return b
}
