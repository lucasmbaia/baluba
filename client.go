package main

import (
	"github.com/lucasmbaia/baluba/gogrpc"
	"golang.org/x/net/context"
	"log"
	"flag"
	"encoding/json"
	"io/ioutil"
	"time"
)

var (
	configDir = flag.String("configDir", "", "")
)

func main() {
	var (
		c     gogrpc.ClientGRPC
		err   error
		dir   []gogrpc.DirectoriesTemplate
		body  []byte
		start time.Time
	)

	start = time.Now()
	flag.Parse()

	if body, err = ioutil.ReadFile(*configDir); err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(body, &dir); err != nil {
		log.Fatal(err)
	}

	if c, err = gogrpc.NewClientGRPC(gogrpc.ClientGRPCConfig{
		Address:	"172.16.95.173:5522",
		ChunkSize:	32768,
		MaxConcurrency:	100,
	}); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if _, err = c.Upload(context.Background(), dir); err != nil {
		log.Fatal(err)
	}

	log.Printf("TEMPO DECORRIDO: %s", time.Since(start))
}
