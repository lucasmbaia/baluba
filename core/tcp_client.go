package core

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"time"
	"net"
	"github.com/lucasmbaia/baluba/core/serializer"
	//"encoding/json"
	//"io/ioutil"
	//"bytes"
)

func UploadFile(name string) error {
	var (
		file	*os.File
		err	error
		conn	net.Conn
		client	= serializer.NewClientSerializer()
		//scanner	*bufio.Scanner
		//buf	= make([]byte, 35 * 1024)
		//n	int
	)

	start := time.Now()
	fmt.Println("PORRA")

	if conn, err = net.Dial("tcp", "192.168.75.129:5522"); err != nil {
		return err
	}

	/*if f, err = json.Marshal(File{
		Name: name,
	}); err != nil {
		return err
	}

	if body, err = encodeGossip(gossip{
		Option:	"teste",
		Body:	f,
	}); err != nil {
		return err
	}

	fmt.Println(string(body))
	conn.Write(body)*/

	if file, err = os.OpenFile(name, os.O_RDONLY, os.ModePerm); err != nil {
		//if file, err = os.Open(name); err != nil {
		return err
	}
	defer func() {
		file.Close()
	}()

	r := bufio.NewReader(file)

	for {
		//var buf = bytes.NewBuffer(make([]byte, 0, 35 * 1024))
		//var body []byte
		var buf = make([]byte, 1024)
		if _, err = r.Read(buf); err != nil {
			if err == io.EOF {
				fmt.Println("DEU OEF")
				break
			}

			fmt.Println(err)
		}

		client.Serializer()
		/*if body, err = encodeGossip(gossip{
			Option:	"file",
			Body:	buf,
		}); err != nil {
			break
		}*/

		conn.Write(buf)
		buf = nil
		//body = nil
	}

	r.Reset(file)
	fmt.Printf("TEMPO DECORRIDO: %s", time.Since(start))
	/*scanner = bufio.NewScanner(file)
	scanner.Buffer(buf, 100 * 1024 * 1024)

	fmt.Println(scanner)
	for scanner.Scan() {
	}

	fmt.Println(file)*/
	return nil
}

