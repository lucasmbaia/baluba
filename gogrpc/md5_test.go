package gogrpc

import (
	"testing"
	"fmt"
)

func Test_CalcMD5(t *testing.T) {
	if h, err := CalcMD5("./md5.go"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(h)
	}
}
