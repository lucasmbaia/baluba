package tcp

import (
	"testing"
)

func Test_StartClient(t *testing.T) {
	StartClient()
}

func Test_SendFiles(t *testing.T) {
	var (
		c   *Client
		err error
		done = make(chan struct{})
	)

	if c, err = NewClient("localhost:5522"); err != nil {
		t.Fatal(err)
	}

	if err = c.SendFiles([]DirectoriesTemplate{
		{Path: "/root/workspace/go/src/github.com/lucasmbaia/baluba/tcp"},
		{Path: "/root/workspace/go/src/github.com/lucasmbaia/baluba/core", Files: []Files{
			{"grpc_server.go"},
		}},
	}); err != nil {
		t.Fatal(err)
	}

	<-done
}
