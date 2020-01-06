package gogrpc

import (
	"encoding/gob"
	"bytes"
	"time"
)

const (
	StatusUnknown	int = 0
	StatusOK	int = 1
	StatusFailed	int = 2
)

type DirectoriesTemplate Directories

type Directories struct {
	Path  string	`json:",omitempty"`
	Files []Files	`json:",omitempty"`
}

type Files struct {
	Name  string  `json:",omitempty"`
}

type Response struct {
	Message	string
	Code	int
}

type Transfer struct {
	Action	  string  `json:",omitempty"`
	Hostname  string  `json:",omitempty"`
	Directory string  `json:",omitempty"`
	FileName  string  `json:",omitempty"`
	Size	  int64	  `json:",omitempty"`
	Content	  []byte  `json:",omitempty"`
}

type Stats struct {
	StartedAt  time.Time
	FinishedAt time.Time
}

func (t *Transfer) Serialize() (b []byte, err error) {
	var (
		result	bytes.Buffer
		encoder	*gob.Encoder
	)

	encoder = gob.NewEncoder(&result)
	if err = encoder.Encode(t); err != nil {
		return
	}

	b = result.Bytes()
	return
}

func Deserialize(b []byte) (tr *Transfer, err error) {
	var (
		decoder	*gob.Decoder
		t	Transfer
	)

	decoder = gob.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&t); err != nil {
		return
	}

	tr = &t
	return
}

func (r *Response) Serialize() (b []byte, err error) {
	var (
		result	bytes.Buffer
		encoder	*gob.Encoder
	)

	encoder = gob.NewEncoder(&result)
	if err = encoder.Encode(r); err != nil {
		return
	}

	b = result.Bytes()
	return
}

func DeserializeResponse(b []byte) (resp *Response, err error) {
	var (
		decoder	*gob.Decoder
		r	Response
	)

	decoder = gob.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&r); err != nil {
		return
	}

	resp = &r
	return
}
