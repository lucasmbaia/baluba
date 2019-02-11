package core

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	//"strconv"
	"time"
	"github.com/lucasmbaia/baluba/core/serializer"
	//"encoding/json"
	//"io/ioutil"
	//"bytes"
)

/****
Cada block de dados transferido vai conter
	Version - 4 Bytes
	Hostname - 4 Bytes
	Option - 4 bytes
	File Name - 4 bytes
	File - 4 bytes

	var pepeca = []byte("tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile_tile")
	var bytesOption = strconv.FormatInt(int64(len(pepeca)/2), 16)

	var b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(pepeca)))

	fmt.Println(bytesOption)
	fmt.Println(hex.EncodeToString(b))

	fmt.Println(binary.LittleEndian.Uint32(b[:]))
***/

func UploadFile(name string) error {
	var (
		file	  *os.File
		err	  error
		//n	  int
		buffer	  = make([]byte, 2048)
		hostname  string
		conn net.Conn
		client = serializer.NewClientSerializer()
		//scanner	*bufio.Scanner
		//buf	= make([]byte, 35 * 1024)
		//n	int
	)

	start := time.Now()
	fmt.Println("PORRA")

	if conn, err = net.Dial("tcp", "192.168.75.129:5522"); err != nil {
		return err
	}


	if hostname, err = os.Hostname(); err != nil {
		return err
	}

	buffer = append(buffer, ConvertUnsigned4Bytes(1)...)
	buffer = append(buffer, []byte("1")...)
	buffer = append(buffer, ConvertUnsigned4Bytes(uint32(len(hostname)))...)
	buffer = append(buffer, []byte(hostname)...)
	buffer = append(buffer, ConvertUnsigned4Bytes(uint32(len("file_name")))...)
	buffer = append(buffer, []byte("file_name")...)
	buffer = append(buffer, ConvertUnsigned4Bytes(uint32(len(name)))...)
	buffer = append(buffer, []byte(name)...)

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

		/*buffer = append(buffer, []byte(bytesOption)...)
		buffer = append(buffer, buf...)*/
		client.Serializer()
		/*if body, err = encodeGossip(gossip{
			Option:	"file",
			Body:	buf,
		}); err != nil {
			break
		}*/

		//conn.Write(append(buffer, append(ConvertUnsigned4Bytes(uint32(n)), buf...)...))
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
