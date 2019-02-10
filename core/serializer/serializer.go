package serializer

import (
	"fmt"
	"github.com/lucasmbaia/baluba/core/serializer/gossip"
	flatbuffers "github.com/google/flatbuffers/go"
)

type Client struct {
	ft  *flatbuffers.Builder
}

func NewClientSerializer() *Client {
	return &Client{
		ft:  flatbuffers.NewBuilder(0),
	}
}

func (c *Client) Serializer() []byte {
	var (
		fo = &gossip.FileObj{Name: "lucas", Size: int64(29)}
	)

	//b := flatbuffers.NewBuilder(0)
	c.ft.Reset()

	name_position := c.ft.CreateString(fo.Name)

	gossip.FileStart(c.ft)
	gossip.FileAddName(c.ft, name_position)
	gossip.FileAddSize(c.ft, fo.Size)

	end := gossip.FileEnd(c.ft)

	c.ft.Finish(end)

	return c.ft.Bytes[c.ft.Head():]
}

func (c *Client) Deserializer(buffer []byte) {
	file := gossip.GetRootAsFile(buffer, 0)

	fmt.Println(string(file.Name()))
	fmt.Println(file.Size())
}
