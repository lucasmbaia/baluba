package gogrpc

import (
	"github.com/kalafut/imohash"
	"fmt"
)

func CalcMD5(file string) (hash string, err error) {
	var md [imohash.Size]byte

	if md, err = imohash.SumFile(file); err != nil {
		return
	}

	hash = fmt.Sprintf("%x", md)
	return
}
