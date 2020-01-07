package gogrpc

import (
	"testing"
	"fmt"
	"context"
	"os"
)

func Test_Open(t *testing.T) {
	if _, err := Open("mysql", DatabaseConfig{
		Username: "nemesis",
		Password: "bg6!OmxR#^QwWu2b",
		Host:	  "gateway-db.dev.nuvem-intera.local",
		Port:	  3306,
		Database: "nemesis_prod",
		Timeout:  "30000ms",
	}); err != nil {
		t.Fatal(err)
	}
}

func Test_ListMysqlDatabases(t *testing.T) {
	var (
		d     Database
		err   error
		data  []string
	)

	/*if d, err = Open("mysql", "nemesis:bg6!OmxR#^QwWu2b@tcp(gateway-db.dev.nuvem-intera.local:3306)/nemesis_prod?parseTime=true&timeout=30000ms"); err != nil {
		t.Fatal(err)
	}*/

	if d, err = Open("mysql", DatabaseConfig{
		Username: "nemesis",
		Password: "bg6!OmxR#^QwWu2b",
		Host:	  "gateway-db.dev.nuvem-intera.local",
		Port:	  3306,
		Database: "nemesis_prod",
		Timeout:  "30000ms",
	}); err != nil {
		t.Fatal(err)
	}

	if data, err = d.ListMysqlDatabases(); err != nil {
		t.Fatal(err)
	}

	fmt.Println(data)
}

func Test_DumpMysqlDatabase(t *testing.T) {
	var (
		body  = make(chan []byte)
		err   error
		d     Database
		file  *os.File
		done  = make(chan struct{})
	)

	if d, err = Open("mysql", DatabaseConfig{
		Username: "nemesis",
		Password: "bg6!OmxR#^QwWu2b",
		Host:	  "gateway-db.dev.nuvem-intera.local",
		Port:	  3306,
		Database: "nemesis_prod",
		Timeout:  "30000ms",
	}); err != nil {
		t.Fatal(err)
	}

	if file, err = os.Create("./teste.sql"); err != nil {
		t.Fatal(err)
	}

	go func() {
		if err = d.DumpMysqlDatabase(context.Background(), "nemesis_copy", body); err != nil {
			t.Fatal(err)
		}

		done <- struct{}{}
	}()

	for {
		select {
		case b := <-body:
			file.Write(b)
		case _ = <-done:
			file.Close()
			return
		}
	}
}
