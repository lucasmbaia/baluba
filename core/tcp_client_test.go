package core

import (
	"fmt"
	"testing"
)

func Test_UploadFile(t *testing.T) {
	//fmt.Println(UploadFile("/root/narutoPROJECT_-_Boruto_013-previa_MQ.mp4"))
	//fmt.Println(UploadFile("/root/Vikings.S05E09.A.Simple.Story.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb.mkv"))
	fmt.Println(UploadFile("/root/go1.13.1.linux-amd64.tar.gz"))
	//UploadFile("/root/workspace/go/src/github.com/lucasmbaia/backup/core/tcp_client.go")
	//UploadFile2("/root/Vikings.S05E09.A.Simple.Story.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb.mkv")
}

func Test_StartClient(t *testing.T) {
	fmt.Println(StartClient([]Directories{
		{Path: "/root", Files: []Files{{Name: "Vikings.S05E09.A.Simple.Story.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb.mkv"}, {Name: "narutoPROJECT_-_Boruto_013-previa_MQ.mp4"}}},
	}))
}
