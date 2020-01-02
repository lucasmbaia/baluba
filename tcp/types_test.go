package tcp

import (
	"testing"
	"fmt"
)

func Test_TSD(t *testing.T) {
	var (
		tr  Transfer
		b   []byte
		err error
		ttr *Transfer
	)

	tr = Transfer{
		Action:	    "action",
		Hostname:   "hostname",
		Directory:  "directory",
	}

	if b, err = tr.Serialize(); err != nil {
		t.Fatal(err)
	}

	if ttr, err = Deserialize(b); err != nil {
		t.Fatal(err)
	}

	fmt.Println(ttr)
}
