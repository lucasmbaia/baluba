package core

import (
	"testing"
	"context"
)

func Test_ClientGRPC(t *testing.T) {
	c, err := NewClientGRPC(ClientGRPCConfig{
		Address: "192.168.75.133:5522",
		//ChunkSize:  (1 << 12),
		ChunkSize:  32768,
		//Compress:   true,
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.Upload(
		context.Background(),
		&[]Directories{
			{Path: "/root", Files: []Files{
				{Name: "go1.13.1.linux-amd64.tar.gz"},
				{Name: "erlang-19.3.6.13-1.el7.centos.x86_64.rpm"},
				//{Name: "Vikings.S05E09.A.Simple.Story.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb.mkv"},
				//{Name: "narutoPROJECT_-_Boruto_013-previa_MQ.mp4"},
				/*{Name: "CentOS-7-x86_64-Minimal-1804.iso"},
				{Name: "ubuntu-mate-16.04.2-desktop-amd64.iso"},
				{Name: "CentOS-7-x86_64-Minimal-1611.iso"},
				{Name: "ubuntu-16.04.1-server-amd64.iso"},
				{Name: "ubuntu-16.04.2-server-amd64.iso"},
				{Name: "ubuntu-18.04-live-server-amd64.iso"},*/
			}},
	}); err != nil {
			t.Fatal(err)
		} else {
			defer c.Close()
		}
	}
