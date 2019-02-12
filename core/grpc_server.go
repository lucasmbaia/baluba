package core

import (
	"net"
	"strconv"
	"fmt"

	"github.com/lucasmbaia/baluba/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type ServerGRPC struct {
	server	*grpc.Server
	port	int
}

type ServerGRPCConfig struct {
	Port  int
}

func NewServerGRPC(cfg ServerGRPCConfig) (s ServerGRPC, err error) {
	if cfg.Port == 0 {
		return s, errors.Errorf("Port must be specified")
	}

	s.port = cfg.Port

	return s, nil
}

func (s *ServerGRPC) Listen() error {
	var (
		l	  net.Listener
		grpcOpts  = []grpc.ServerOption{}
		err	  error
	)

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
	return nil
}

func (s *ServerGRPC) Close() {
	if s.server != nil {
		s.server.Stop()
	}

	return
}
