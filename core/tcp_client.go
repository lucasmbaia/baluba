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
	"path/filepath"
	"strings"
	//"internal/poll"
	//"runtime"
	//"syscall"
	//"compress/gzip"
	//"encoding/json"
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
	src	  *os.File
	done	  bool
	size	  int64
	position  int64
	client	  *serializer.Client
}

func NewReader(src *os.File, size int64, client *serializer.Client) *Reader {
	return &Reader{
		src:	src,
		size:	size,
		client:	client,
	}
}

func (r *Reader) Read(b []byte) (n int, err error) {
	if r.done {
		return 0, io.EOF
	}

	if r.size - r.position < int64(len(b)) {
		r.done = true
	}

	//var buffer = make([]byte, len(b) - len(r.client.SerializerGossip(gossip.GossipObj{Option: "file_name"})))
	r.src.Seek(r.position, 0)
	n, err = r.src.Read(b)
	r.position += int64(n)

	//copy(b, r.client.SerializerGossip(gossip.GossipObj{Option: "file_name", Body: buffer}))
	//r.client.DeserializerGossip(b)

	return n, err
}

func StartClient(directories []Directories) error {
	var (
		err	error
		conn     net.Conn
	)

	checkDirectories := func(directories []Directories, path string) (int, bool) {
		for idx, d := range directories {
			if d.Path == path {
				return idx, true
			}
		}

		return 0, false
	}

	for idx, d := range directories {
		if _, err = os.Stat(d.Path); os.IsNotExist(err) || err != nil {
			return err
		}

		if len(d.Files) > 0 {
			for _, file := range d.Files {
				if _, err = os.Stat(fmt.Sprintf("%s/%s", d.Path, file.Name)); os.IsNotExist(err) || err != nil {
					return err
				}
			}
		} else {
			if err = filepath.Walk(d.Path, func(path string, info os.FileInfo, e error) error {
				var (
					file	string
					index	int
					exists	bool
					dir	[]string
				)


				if info.IsDir() {
					if _, exists = checkDirectories(directories, path); !exists {
						directories = append(directories, Directories{
							Path:	path,
						})
					}

					return nil
				}

				file = strings.Replace(path, fmt.Sprintf("%s/", d.Path), "", 1)
				if len(strings.Split(file, "/")) == 1 && !info.IsDir() {
					directories[idx].Files = append(directories[idx].Files, Files{Name: file})
				} else {
					dir = strings.Split(path, "/")

					if index, exists = checkDirectories(directories, strings.Join(dir[:len(dir) -1], "/")); exists {
						directories[index].Files = append(directories[index].Files, Files{Name: dir[len(dir)-1]})
					} else {
						directories = append(directories, Directories{
							Path:	strings.Join(dir[:len(dir) -1], "/"),
							Files:	[]Files{{Name: dir[len(dir)-1]}},
						})
					}
				}

				return nil
			}); err != nil {
				return err
			}
		}
	}

	if conn, err = net.Dial("tcp", "172.16.95.171:5522"); err != nil {
		return err
	}

	for _, d := range directories {
		for _, f := range d.Files {
			//if err = FastUploadFile(fmt.Sprintf("%s/%s", d.Path, f.Name), conn); err != nil {
			if err = UploadFileBuffering(fmt.Sprintf("%s/%s", d.Path, f.Name), conn); err != nil {
				return err
			}
		}
	}

	return nil
}

func FastUploadFile(name string, conn net.Conn) error {
	var (
		err	error
		file	*os.File
		stat	os.FileInfo
		r	*Reader
		w	*bufio.Writer
		client	= serializer.NewClientSerializer()
		written	int64
	)

	if file, err = os.OpenFile(name, os.O_RDONLY, os.ModePerm); err != nil {
		return err
	}
	defer file.Close()

	if stat, err = file.Stat(); err != nil {
		return err
	}

	r = NewReader(file, stat.Size(), client)
	w = bufio.NewWriter(conn)

	if written, err = io.CopyN(w, r, stat.Size()); err != nil {
		return err
	}

	fmt.Println(stat.Size())
	fmt.Println(written)

	return nil
}

func UploadFileBuffering(name string, conn net.Conn) error {
	var (
		buffer	  []byte
		hostname  string
		err	  error
		r	  *bufio.Reader
		client	  = serializer.NewClientSerializer()
		n	  int
		file	  *os.File
	)

	if hostname, err = os.Hostname(); err != nil {
		return err
	}

	buffer = createHeader("transfer", hostname, name)
	fmt.Println(buffer)
	fmt.Println(n, client)

	if file, err = os.OpenFile(name, os.O_RDONLY, os.ModePerm); err != nil {
		return err
	}
	defer file.Close()

	r = bufio.NewReader(file)

	for {
		var buf = make([]byte, 32768)
		if _, err = r.Read(buf); err != nil {
			if err == io.EOF {
				break
			}
		}

		//conn.Write(append(buffer, append(ConvertUnsigned4Bytes(uint32(n)), buf...)...))
		//fmt.Println(conn.Write(client.SerializerGossip(gossip.GossipObj{Option: "transfer", Body: buf})))
		fmt.Println(conn.Write(buf))

		time.Sleep(2 * time.Second)
	}

	return nil
}

func createHeader(option, hostname, file string) []byte {
	var buffer []byte

	buffer = append(buffer, ConvertUnsigned4Bytes(1)...)
	buffer = append(buffer, []byte("1")...)
	buffer = append(buffer, ConvertUnsigned4Bytes(uint32(len(hostname)))...)
	buffer = append(buffer, []byte(hostname)...)
	buffer = append(buffer, ConvertUnsigned4Bytes(uint32(len(option)))...)
	buffer = append(buffer, []byte(option)...)
	buffer = append(buffer, ConvertUnsigned4Bytes(uint32(len(file)))...)
	buffer = append(buffer, []byte(file)...)

	return buffer
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
	r := NewReader(file, stat.Size(), client)
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
