package gogrpc

import (
	"testing"
	"fmt"
)

func Test_ListFiles(t *testing.T) {
	var (
		d   []Directories
		err error
	)

	if d, err = ListFiles([]DirectoriesTemplate{
		{Path: "/root/workspace/go/src/github.com/lucasmbaia/baluba/core"},
		{Path: "/root/workspace/go/src/github.com/lucasmbaia/baluba/tcp"},
		/*{Path: "/root/workspace/go/src/github.com/lucasmbaia/baluba/core", Files: []Files{
			{"grpc_server.go"},
		}},*/
	}, true, 0); err != nil {
		t.Fatal(err)
	}

	fmt.Println(d)
}
