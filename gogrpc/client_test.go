package gogrpc

import (
	"testing"
	"context"
)

func Test_ClientGRPC(t *testing.T) {
	c, err := NewClientGRPC(ClientGRPCConfig{
		Address:    "192.168.75.133:5522",
		ChunkSize:  32768,
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.Upload(
		context.Background(),
		[]DirectoriesTemplate{
			{Path: "/root/workspace/go/src/github.com/lucasmbaia/baluba/tcp"},
			{Path: "/root/workspace/go/src/github.com/lucasmbaia/baluba/core", Files: []Files{
				{"grpc_server.go"},
			}},
		},
	); err != nil {
		t.Fatal(err)
	} else {
		defer c.Close()
	}
}

