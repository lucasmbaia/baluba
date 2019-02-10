package core

import (
	"fmt"
	"testing"
)

func Test_UploadFile(t *testing.T) {
	fmt.Println(UploadFile("/root/Vikings.S05E20.Ragnarok.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb.mkv"))
	//UploadFile("/root/workspace/go/src/github.com/lucasmbaia/backup/core/tcp_client.go")
	//UploadFile2("/root/Vikings.S05E09.A.Simple.Story.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb.mkv")
}
