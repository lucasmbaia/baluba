package tcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

const (
	rootPath  = "/baluba"
)

type connection struct {
	sync.RWMutex

	write *bufio.Writer
	read  *bufio.Reader
	conn  net.Conn
	files map[string]*files
	host  string
}

type files struct {
	content	[]byte
	size	int64
	writer	int64
	f	*os.File
	done	chan struct{}
}

func NewServerTcp(ctx context.Context) {
	var (
		err error
		l   net.Listener
	)

	if l, err = net.Listen("tcp", ":5522"); err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	go func() {
		for {
			var c = &connection{}

			if c.conn, err = l.Accept(); err != nil {
				log.Fatalf(fmt.Sprintf("Error to accept connection: %s\n", err.Error()))
				continue
			}

			c.write = bufio.NewWriter(c.conn)
			c.read = bufio.NewReader(c.conn)
			c.files = make(map[string]*files)

			go c.handleConnection()
		}
	}()

	<-ctx.Done()
}

func (c *connection) handleConnection() {
	defer func() {
		if addr, ok := c.conn.RemoteAddr().(*net.TCPAddr); ok {
			log.Printf("Connection Close if node IP %s", addr.IP.String())
		}
		c.conn.Close()
	}()

	var (
		err	error
		buffer	= make([]byte, 32768)
		n	int
	)

	for {
		if n, err = c.conn.Read(buffer); err != nil {
			if err != io.EOF {
				log.Printf("Error to read bytes: %s\n", err.Error())
			}

			return
		}

		//c.decision(buffer[:n])
		fmt.Println("TOMA", n)
	}
}

func (c *connection) decision(buffer []byte) {
	var (
		t   *Transfer
		err error
	)

	if t, err = Deserialize(buffer); err != nil {
		fmt.Println(err)
	}

	fmt.Println(t)
	switch t.Action{
	case "create_file":
		if err = c.createFile(t); err != nil {
			fmt.Println("DEU MERDA: ", t.FileName)
			c.response(Response{
				Message:  err.Error(),
				Code:	  StatusFailed,
			})
		} else {
			fmt.Println("STATUS OK: ", t.FileName)
			c.response(Response{
				Code: StatusOK,
			})
		}
	case "transfer_file":
	}
}

func (c *connection) response(r Response) (err error) {
	var (
		b []byte
	)

	if b, err = r.Serialize(); err != nil {
		return
	}

	fmt.Println("WRITE")
	_, err = c.conn.Write(b)
	fmt.Println("RETURN")
	return
}

func (c *connection) createFile(t *Transfer) (err error) {
	var (
		file	  *os.File
		fullPath  string
		fullFile  string
	)

	fullPath = fmt.Sprintf("%s/%s/%s", rootPath, t.Hostname, t.Directory)
	if _, err = os.Stat(fullPath); os.IsNotExist(err) {
		if err = os.MkdirAll(fullPath, os.ModePerm); err != nil {
			return err
		}
	}

	fullFile = fmt.Sprintf("%s/%s", fullPath, t.FileName)
	if file, err = os.Create(fullFile); err != nil {
		return err
	}

	c.Lock()
	c.files[fullFile] = &files{
		size:   t.Size,
		f:      file,
		writer: 0,
	}
	c.Unlock()

	fmt.Println("DEU UNLOCK")
	return
}
