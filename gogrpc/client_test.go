package gogrpc

import (
	"testing"
	"context"
)

func Test_ClientGRPC(t *testing.T) {
	c, err := NewClientGRPC(ClientGRPCConfig{
		Address:	"192.168.75.133:5522",
		ChunkSize:	32768,
		MaxConcurrency:	3000,
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.Upload(
		context.Background(),
		[]DirectoriesTemplate{
			/*{Path: "/root/teste-baluba", Files: []Files{
				{"ubuntu-mate-16.04.2-desktop-amd64.iso"},
			}},*/
			//{Path: "/root/workspace/go/src/github.com/lucasmbaia/baluba/tcp"},
			{Path: "/root/teste-baluba/small"},
			//{Path: "/root/workspace/go/src/github.com/lucasmbaia/baluba/core", Files: []Files{
			//	{"grpc_server.go"},
			//}},
		},
	); err != nil {
		t.Fatal(err)
	} else {
		defer c.Close()
	}
}

func Test_UploadDatabases(t *testing.T) {
	c, err := NewClientGRPC(ClientGRPCConfig{
		Address:	"192.168.207.128:5522",
		ChunkSize:	32768,
		MaxConcurrency:	3000,
		DatabaseConfig:	&DatabaseConfig{
			Username: "root",
			Password: "123456",
			Host:     "127.0.0.1",
			Port:     3306,
			Timeout:  "30000ms",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.UploadDatabases(context.Background()); err != nil {
		t.Fatal(err)
	} else {
		defer c.Close()
	}
}

