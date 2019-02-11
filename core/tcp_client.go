package core

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	//"strconv"
	"github.com/lucasmbaia/baluba/core/serializer"
	//"github.com/lucasmbaia/baluba/core/serializer/gossip"
	"time"
	//"internal/poll"
	//"runtime"
	//"syscall"
	//"compress/gzip"
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

type Reader struct {
	//src	io.Reader
	src	*os.File
	done 	bool
	size	int64
	position	int64
	client		*serializer.Client
	//pfd	poll.FD
}

func NewReader(src *os.File, fd uintptr, size int64, client *serializer.Client) *Reader {
	return &Reader{
		src:	src,
		size:	size,
		client:	client,
		/*pfd:	poll.FD{
			Sysfd:		int(fd),
			IsStream:	true,
			ZeroReadIsEOF:	true,
		},*/
	}
}

func (r *Reader) Read(b []byte) (n int, err error) {
	if r.done {
		return 0, io.EOF
	}

	//var off int64 = 1024
	//var buffer []byte
	if r.size - r.position < 1024 {
		r.done = true
		//off = r.size - r.position
	}

	r.src.Seek(r.position, 0)
	n, err = r.src.Read(b)
	r.position += int64(n)

	fmt.Println(b)
	//b = r.client.SerializerGossip(gossip.GossipObj{Option: "file_name", Body: buffer})
	return n, err
	/*n, err = r.pfd.Read(b)
	runtime.KeepAlive(r)
	return n, err*/


	/*if r.done {
		return 0, io.EOF
	}

		p = []byte("pepca")

	r.done = true
	return 1, nil*/
}

func UploadFile(name string) error {
	var (
		file     *os.File
		err      error
		//n        int
		buffer   []byte
		hostname string
		conn     net.Conn
		//conn     *net.TCPConn
		send	= make(chan []byte)
		client   = serializer.NewClientSerializer()
		//scanner	*bufio.Scanner
		//buf	= make([]byte, 35 * 1024)
		//n	int
	)

	start := time.Now()
	fmt.Println("PORRA")
	fmt.Println(conn)

	/*tcpAddr, err := net.ResolveTCPAddr("tcp", "172.16.95.171:5522")
	if err != nil {
		return err
	}

	if conn, err = net.DialTCP("tcp", nil, tcpAddr); err != nil {
		return err
	}*/
	if conn, err = net.Dial("tcp", "172.16.95.171:5522"); err != nil {
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


	go func() {
		for {
			select {
			case b := <-send:
				conn.Write(b)
			}
		}
	}()
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
		return err
	}
	defer func() {
		file.Close()
	}()

	//r := bufio.NewReader(file)
	stat, _ := file.Stat()
	r := NewReader(file, file.Fd(), stat.Size(), client)
	bw := bufio.NewWriter(conn)



	if _, err = io.CopyN(bw, r, stat.Size()); err != nil {
		return err
	}

	/*for {
		var buf = make([]byte, 1024)
		if _, err = r.Read(buf); err != nil {
			if err == io.EOF {
				fmt.Println("DEU OEF")
				break
			}

			fmt.Println(err)
		}

		//client.SerializerGossip(gossip.GossipObj{Option: "file_name", Body: buf})
		//client.SerializerGossip(gossip.GossipObj{Option: "file_name", Body: buf})
		//fmt.Println(len(append(buffer, append(ConvertUnsigned4Bytes(uint32(n)), buf...)...)))
		//conn.Write(append(buffer, append(ConvertUnsigned4Bytes(uint32(n)), buf...)...))
		//conn.Write(client.SerializerGossip(gossip.GossipObj{Option: "file_name", Body: buf}))
		conn.Write(buf)
		buf = nil
		//body = nil
	}*/

	//r.Reset(file)
	fmt.Printf("TEMPO DECORRIDO: %s", time.Since(start))
	/*scanner = bufio.NewScanner(file)
	scanner.Buffer(buf, 100 * 1024 * 1024)

	fmt.Println(scanner)
	for scanner.Scan() {
	}

	fmt.Println(file)*/
	return nil
}
