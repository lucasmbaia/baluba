package core

import (
	"testing"
)

func Test_ServerGRPC(t *testing.T) {
	s, err := NewServerGRPC(ServerGRPCConfig{
		Port:	  5522,
		RootPath: "/root",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Listen(); err != nil {
		t.Fatal(err)
	} else {
		defer s.Close()
	}
}
