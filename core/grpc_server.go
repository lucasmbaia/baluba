package core

import (
	"net"
	"strconv"
	"fmt"
	"io"
	"sync"
	"os"

	"google.golang.org/grpc/metadata"
	"github.com/lucasmbaia/baluba/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type ServerGRPC struct {
	sync.RWMutex

	server	  *grpc.Server
	port	  int
	files	  map[string]*files
	response  chan response
	rootPath  string
	//hosts	  map[string]Host
}

type response struct {
	stream	baluba.BalubaService_UploadServer
	file	string
}

type files struct {
	content	[]byte
	size	int64
	writer	int64
	f	*os.File
	done	chan struct{}
	stream	baluba.BalubaService_UploadServer
}

type ServerGRPCConfig struct {
	Port	  int
	RootPath  string
}

func NewServerGRPC(cfg ServerGRPCConfig) (s ServerGRPC, err error) {
	if cfg.Port == 0 {
		return s, errors.Errorf("Port must be specified")
	}

	if cfg.RootPath == "" {
		return s, errors.Errorf("Root Path must be specified")
	}

	s.port = cfg.Port
	s.files = make(map[string]*files)
	s.response = make(chan response)
	//s.hosts = make(map[string]Host)
	s.rootPath = cfg.RootPath

	return s, nil
}

func (s *ServerGRPC) Listen() error {
	var (
		l	  net.Listener
		grpcOpts  = []grpc.ServerOption{}
		err	  error
	)

	fmt.Println(fmt.Sprintf(":%s", strconv.Itoa(s.port)))
	if l, err = net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(s.port))); err != nil {
		return err
	}

	s.server = grpc.NewServer(grpcOpts...)
	baluba.RegisterBalubaServiceServer(s.server, s)
	if err = s.server.Serve(l); err != nil {
		return err
	}

	return nil
}

func (s *ServerGRPC) Upload(stream baluba.BalubaService_UploadServer) error {
	var (
		chunk *baluba.Chunk
		err   error
		//done  = make(chan struct{})
	)

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		fmt.Println(md)
	}
	//fmt.Println(stream.Context())
	//s.files["teste"] = make(chan []byte)

	//go s.write(done, s.files["teste"])

	SEND := func(stream baluba.BalubaService_UploadServer, message, f string, code baluba.UploadStatusCode) {
		fmt.Println("PORRA")
		if stream != nil {
			fmt.Println(stream)
			if err = stream.Send(&baluba.UploadStatus{
				Message:  message,
				Code:	  code,
				FileName: f,
			}); err != nil {
				fmt.Println(errors.Wrapf(err, "failed to send status code"))
			}
		}
	}

	go func() {
		for {
			select {
			case resp := <-s.response:
				SEND(resp.stream, "", resp.file, baluba.UploadStatusCode_Ok)
			}
		}
	}()

	for {
		if chunk, err = stream.Recv(); err != nil {
			if err == io.EOF {
				return nil
			}

			if err != nil {
				fmt.Println("TOMA NO CU PORRA", err, stream)
				SEND(stream, fmt.Sprintf("%v, failed unexpectadely while reading chunks from stream", err), "", baluba.UploadStatusCode_Failed)
				return err
			}
		}

		//s.route(chunk)
		switch chunk.Action {
		case "create_file":
			if err = s.createFile(chunk, stream); err != nil {
				SEND(stream, err.Error(), chunk.Name, baluba.UploadStatusCode_Failed)
			} else {
				SEND(stream, "", chunk.Name, baluba.UploadStatusCode_Ok)
			}
		case "transfer":
			s.writer(chunk)
		}

		//s.route(chunk)
		//s.files["teste"] <- chunk.Content
	}

	/*END:
	if err = stream.SendAndClose(&baluba.UploadStatus{
		Message: "Upload received with success",
		Code:    baluba.UploadStatusCode_Ok,
	}); err != nil {
		return errors.Wrapf(err, "failed to send status code")
	}*/

	return nil
}

func (s *ServerGRPC) writer(chunk *baluba.Chunk) {
	s.Lock()
	s.files[chunk.Name].f.Write(chunk.Content)
	s.files[chunk.Name].writer += int64(len(chunk.Content))

	if s.files[chunk.Name].writer == s.files[chunk.Name].size {
		fmt.Println("fecha aqui")
		s.files[chunk.Name].f.Close()
		s.response <- response{
			stream: s.files[chunk.Name].stream,
			file:	chunk.Name,
		}
		delete(s.files, chunk.Name)
	}

	s.Unlock()
}

func (s *ServerGRPC) createFile(chunk *baluba.Chunk, stream baluba.BalubaService_UploadServer) error {
	var (
		err	  error
		file	  *os.File
		fullPath  string
	)

	fullPath = fmt.Sprintf("%s/%s/%s", s.rootPath, chunk.Hostname, chunk.Directory)
	if _, err = os.Stat(fullPath); os.IsNotExist(err) {
		if err = os.MkdirAll(fullPath, os.ModePerm); err != nil {
			return err
		}
	}

	if file, err = os.Create(fmt.Sprintf("%s/%s", fullPath, chunk.Name)); err != nil {
		return err
	}

	s.Lock()
	s.files[chunk.Name] = &files{
		size:	chunk.Size,
		f:	file,
		writer:	0,
		stream:	stream,
	}
	s.Unlock()

	return nil
}

func (s *ServerGRPC) route(chunk *baluba.Chunk) {
	var err error

	switch chunk.Action {
	case "create_file":
		var (
			file	  *os.File
			fullPath  string
		)

		fullPath = fmt.Sprintf("%s/%s/%s", s.rootPath, chunk.Hostname, chunk.Directory)
		if _, err = os.Stat(fullPath); os.IsNotExist(err) {
			if err = os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return
			}
		}

		if file, err = os.Create(fmt.Sprintf("%s/%s", fullPath, chunk.Name)); err != nil {
			return
		}

		s.Lock()
		s.files[chunk.Name] = &files{
			size:	chunk.Size,
			f:	file,
			writer:	0,
		}
		//go w.write(s.files[chunk.Name])
		s.Unlock()
	case "transfer":
		go func() {
			s.Lock()
			s.files[chunk.Name].f.Write(chunk.Content)
			s.files[chunk.Name].writer += int64(len(chunk.Content))

			if s.files[chunk.Name].writer == s.files[chunk.Name].size {
				fmt.Println("fecha aqui")
				s.files[chunk.Name].f.Close()
			}
			s.Unlock()
		}()
	}
	/*s.Lock()
	if _, ok := s.files[chunk.Name]; ok {
		s.files[chunk.Name] <- chunk.Content
	} else {
		var done = make(chan struct{})
		s.files[chunk.Name] = make(chan []byte)
		go s.write(done, s.files[chunk.Name])
		s.files[chunk.Name] <- chunk.Content
	}
	s.Unlock()*/
}

func (s *ServerGRPC) write(done chan struct{}, buffer <-chan []byte) {
	for {
		select {
		case _ = <-buffer:
		case _ = <-done:
			return
		}
	}
}

func (s *ServerGRPC) Close() {
	if s.server != nil {
		s.server.Stop()
	}

	return
}
