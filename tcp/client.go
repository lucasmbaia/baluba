package tcp

import (
	"net"
	"fmt"
	"bufio"
	"io"
	"bytes"
	"os"
	"log"
	//"strings"
)

type Client struct {
	write *bufio.Writer
	read  *bufio.Reader
	conn  net.Conn
}

func NewClient(addr string) (c *Client, err error) {
	c = &Client{}

	if c.conn, err = net.Dial("tcp", addr); err != nil {
		return
	}

	c.write = bufio.NewWriter(c.conn)
	c.read = bufio.NewReader(c.conn)

	return
}

func (c *Client) SendFiles(dt []DirectoriesTemplate) (err error) {
	var (
		dir []Directories
	)

	if dir, err = ListFiles(dt); err != nil {
		return
	}

	go c.receiver()

	for _, d := range dir {
		for _, f := range d.Files {
			go c.transferFile(d.Path, f.Name)
		}
	}

	return
}

func (c *Client) receiver() {
	var (
		err	error
		buffer	= make([]byte, 512)
		n	int
	)

	for {
		var resp *Response
		if n, err = c.read.Read(buffer); err != nil {
			if err != io.EOF {
				log.Printf("Error to read bytes: %s\n", err.Error())
			}

			continue
		}

		if resp, err = DeserializeResponse(buffer[:n]); err != nil {
			return
		}

		fmt.Println(resp)
	}
}

func (c *Client) transferFile(path, fileName string) (err error) {
	var (
		b	[]byte
		t	Transfer
		file	*os.File
		stat	os.FileInfo
	)

	if file, err = os.OpenFile(fmt.Sprintf("%s/%s", path, fileName), os.O_RDONLY, os.ModePerm); err != nil {
		return
	}
	defer file.Close()

	if stat, err = file.Stat(); err != nil {
		return
	}

	t = Transfer{
		Action:	    "create_file",
		Hostname:   "lucas",
		Directory:  path,
		FileName:   fileName,
		Size:	    stat.Size(),
	}

	if b, err = t.Serialize(); err != nil {
		return
	}

	if _, err = io.CopyN(c.write, bytes.NewReader(b), int64(len(b))); err != nil {
		return
	}

	return
}

func StartClient() (err error) {
	var c = &connection{}

	if c.conn, err = net.Dial("tcp", "localhost:5522"); err != nil {
		return
	}

	c.write = bufio.NewWriter(c.conn)
	c.read = bufio.NewReader(c.conn)

	var tr = Transfer{
		Action:     "action",
		Hostname:   "hostname",
		Directory:  "directory",
	}

	var b []byte
	if b, err = tr.Serialize(); err != nil {
		return
	}

	var buf = bytes.NewReader(b)
	if _, err = io.CopyN(c.write, buf, int64(len(b))); err != nil {
		fmt.Println(err)
	}

	//s1 := strings.NewReader("lucas")
	//r := bufio.NewReaderSize(s1, 16)

	//_, err = io.CopyN(c.write, s1, 5)
	//fmt.Println(err)
	//c.conn.Write([]byte("lucas"))
	//c.write.Write([]byte("lucas"))
	//fmt.Println(c.write.Write([]byte("lucas\n")))
	//fmt.Fprint(c.write, "lucas")
	//c.write.Flush()

	return
}
