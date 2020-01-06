package main

import (
	"github.com/lucasmbaia/baluba/gogrpc"
	"log"
)

func main() {
	var (
		s   gogrpc.ServerGRPC
		err error
	)

	if s, err = gogrpc.NewServerGRPC(gogrpc.ServerGRPCConfig{
		Port:	  5522,
		RootPath: "/root/baluba",
	}); err != nil {
		log.Fatal(err)
	}

	if err = s.Listen(); err != nil {
		log.Fatal(err)
	}
	defer s.Close()
}
