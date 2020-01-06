package gogrpc

import (
	"net"
	"strconv"
	"fmt"
	"io"
	"sync"
	"os"

	"google.golang.org/grpc/metadata"
	"github.com/lucasmbaia/baluba/proto"
	"golang.org/x/net/context"
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
	hosts	  map[string]*Host
}

type Host struct {
	sync.RWMutex

	stream	  map[string]baluba.BalubaService_UploadServer
	files	  map[string]*files
	infos	  chan *Chunk
}

type Chunk struct {
	chunk	*baluba.Chunk
	stream	baluba.BalubaService_UploadServer
}

type response struct {
	file	string
	path	string
	err	error
}

type files struct {
	content	[]byte
	size	int64
	writer	int64
	f	*os.File
	done	chan struct{}
	stream	baluba.BalubaService_UploadServer
	hash	string
}

type ServerGRPCConfig struct {
	Port	  int
	RootPath  string
}

var rootPath string

func NewServerGRPC(cfg ServerGRPCConfig) (s ServerGRPC, err error) {
	if cfg.Port == 0 {
		return s, errors.Errorf("Port must be specified")
	}

	if cfg.RootPath == "" {
		return s, errors.Errorf("Root Path must be specified")
	}

	s.port = cfg.Port
	s.response = make(chan response)
	s.hosts = make(map[string]*Host)
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

	rootPath = s.rootPath
	s.server = grpc.NewServer(grpcOpts...)
	baluba.RegisterBalubaServiceServer(s.server, s)
	if err = s.server.Serve(l); err != nil {
		return err
	}

	return nil
}

func (s *ServerGRPC) Create(ctx context.Context, bfiles *baluba.Files) (status *baluba.UploadStatus, err error) {
	var (
		hostname  string
		md	  map[string][]string
		ok	  bool
		val	  []string
		fpath	  string
		ffile	  string
		file	  *os.File
	)

	defer func() {
		if err != nil {
			s.Lock()
			if hostname != "" {
				delete(s.hosts, hostname)
			}
			s.Unlock()
		}
	}()

	if md, ok = metadata.FromIncomingContext(ctx); ok {
		if val, ok = md["hostname"]; ok {
			hostname = val[0]
		} else {
			status = &baluba.UploadStatus{
				Message:  "Hostname Unknow",
				Code:	  baluba.UploadStatusCode_Unknown,
			}

			return
		}
	} else {
		status = &baluba.UploadStatus{
			Message:  "Error to get infos from context",
			Code:	  baluba.UploadStatusCode_Failed,
		}

		return
	}

	s.Lock()
	s.hosts[hostname] = &Host{
		files:	  make(map[string]*files),
		infos:	  make(chan *Chunk),
		stream:	  make(map[string]baluba.BalubaService_UploadServer),
	}
	s.Unlock()

	for _, f := range (*bfiles).File {
		fpath = fmt.Sprintf("%s/%s%s", s.rootPath, hostname, f.Directory)
		if _, err = os.Stat(fpath); os.IsNotExist(err) {
			if err = os.MkdirAll(fpath, os.ModePerm); err != nil {
				return
			}
		}

		ffile = fmt.Sprintf("%s/%s", fpath, f.Name)
		if file, err = os.Create(ffile); err != nil {
			return
		}

		s.Lock()
		s.hosts[hostname].files[ffile] = &files{
			size:	f.Size,
			f:	file,
			writer:	0,
			hash:	f.Hash,
		}
		s.Unlock()
	}

	status = &baluba.UploadStatus{
		Code:	  baluba.UploadStatusCode_Ok,
	}

	return
}

func (s *ServerGRPC) Upload(stream baluba.BalubaService_UploadServer) (err error) {
	var (
		chunk	  *baluba.Chunk
		hostname  string
		md	  map[string][]string
		ok	  bool
		val	  []string
	)

	SEND := func(stream baluba.BalubaService_UploadServer, status *baluba.UploadStatus) {
		if stream != nil {
			stream.Send(status)
		}
	}

	if md, ok = metadata.FromIncomingContext(stream.Context()); ok {
		if val, ok = md["hostname"]; ok {
			hostname = val[0]
		} else {
			SEND(stream, &baluba.UploadStatus{
				Message:  "Hostname Unknow",
				Code:	  baluba.UploadStatusCode_Unknown,
			})

			return
		}
	} else {
		SEND(stream, &baluba.UploadStatus{
			Message:  "Hostname Unknow",
			Code:	  baluba.UploadStatusCode_Unknown,
		})

		return
	}

	s.Lock()
	if _, ok = s.hosts[hostname]; ok {
		go s.hosts[hostname].transfer(s.hosts[hostname].infos)
	}
	s.Unlock()

	for {
		if chunk, err = stream.Recv(); err != nil {
			if err == io.EOF {
				return
			}

			if err != nil {
				SEND(stream, &baluba.UploadStatus{
					Message:  fmt.Sprintf("%v, failed unexpectadely while reading chunks from stream", err),
					Code:	  baluba.UploadStatusCode_Failed,
				})
			}

			return
		}

		s.Lock()
		s.hosts[chunk.Hostname].infos <- &Chunk{
			chunk:	  chunk,
			stream:	  stream,
		}
		s.Unlock()
	}

	return
}

func (h *Host) transfer(chunk chan *Chunk) {
	for {
		select {
		case c := <-chunk:
			var (
				ffile string
				err   error
				hash  string
			)

			ffile = fmt.Sprintf("%s/%s%s/%s", rootPath, c.chunk.Hostname, c.chunk.Directory, c.chunk.Name)

			h.Lock()
			if _, err = h.files[ffile].f.Write(c.chunk.Content); err != nil {
				c.stream.Send(&baluba.UploadStatus{
					Code:	  baluba.UploadStatusCode_Failed,
					Message:  err.Error(),
				})

				h.files[ffile].f.Close()
				delete(h.files, c.chunk.Name)
			}
			h.files[ffile].writer += int64(len(c.chunk.Content))

			if h.files[ffile].writer == h.files[ffile].size {
				h.files[ffile].f.Close()

				if hash, err = CalcMD5(ffile); err != nil {
					c.stream.Send(&baluba.UploadStatus{
						Code:	  baluba.UploadStatusCode_Failed,
						Message:  err.Error(),
					})
				}

				if h.files[ffile].hash == hash {
					c.stream.Send(&baluba.UploadStatus{
						Code: baluba.UploadStatusCode_Ok,
					})
				} else {
					c.stream.Send(&baluba.UploadStatus{
						Code:	  baluba.UploadStatusCode_Failed,
						Message:  "Hashs of files are different",
					})
				}

				delete(h.files, c.chunk.Name)
			}
			h.Unlock()
		}
	}
}

func (s *ServerGRPC) Close() {
	if s.server != nil {
		s.server.Stop()
	}

	return
}
