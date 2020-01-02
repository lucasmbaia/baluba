package core

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"github.com/lucasmbaia/baluba/core/serializer"
	//"os"
)

const (
	defaultPort = ":5522"
)

type connection struct {
	write *bufio.Writer
	read  *bufio.Reader
	conn  net.Conn
}

type gossipOld struct {
	Option string
	Body   []byte
	Error  error
}

type File struct {
	Name string
	Size int64
}

func decodeGossip(b []byte) (gossipOld, error) {
	var (
		g   gossipOld
		err error
	)

	err = json.Unmarshal(b, &g)
	return g, err
}

func encodeGossip(g gossipOld) ([]byte, error) {
	var (
		buf bytes.Buffer
		err error
	)

	var encoder *gob.Encoder = gob.NewEncoder(&buf)
	err = encoder.Encode(g)

	return buf.Bytes(), err
}

func NewServerTcp(ctx context.Context) {
	var (
		err error
		l   net.Listener
	)

	if l, err = net.Listen("tcp", defaultPort); err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	go func() {
		for {
			var c = &connection{}

			if c.conn, err = l.Accept(); err != nil {
				log.Fatalf(fmt.Sprintf("Error to accept connection: %s\n", err.Error()))
			}

			c.write = bufio.NewWriter(c.conn)
			c.read = bufio.NewReader(c.conn)

			go c.handleConnection()
		}
	}()

	<-ctx.Done()
}

func (c *connection) handleConnection() {
	//var file *os.File

	defer func() {
		if addr, ok := c.conn.RemoteAddr().(*net.TCPAddr); ok {
			log.Printf("Connection Close if node IP %s", addr.IP.String())
		}
		c.conn.Close()
		//file.Close()
	}()

	var (
		err    error
		buffer = make([]byte, 32768)
		client	= serializer.NewClientSerializer()
		n      int
		//total  int64
	)

	fmt.Println(client)
	/*if file, err = os.Create("/root/Vikings.S05E09.A.Simple.Story.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb.mkv"); err != nil {
		return
	}*/

	//total = 0

	for {
		//if n, err = c.read.Read(buffer); err != nil {
		if n, err = io.ReadFull(c.read, buffer); err != nil {
			if err != io.EOF {
				log.Printf("Error to read bytes: %s\n", err.Error())
			}
			return
		}

		fmt.Println(n)
		//fmt.Println(c.read.Peek(n))
		//fmt.Println(c.read.Size())
		//client.DeserializerGossip(buffer[:n])
		/*var g gossipOld

		if g, err = decodeGossip(buffer[:n]); err != nil {
			fmt.Println("DEU MERDA", err)
			continue
		}*/

		/*switch g.Option {
		case "new_file":
			var f File

			if err = json.Unmarshal(g.Body, &f); err != nil {
				continue
			}
		}*/
		//total += int64(n)
		//fmt.Println(total)
		//file.Write(buffer[:n])
		//fmt.Println(g)
	}
}
