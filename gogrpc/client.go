package gogrpc

import (
	"io"
	"os"
	//"strings"
	//"path/filepath"
	"fmt"
	//"sync"

	"github.com/lucasmbaia/baluba/proto"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	_ "google.golang.org/grpc/encoding/gzip"
)

type ClientGRPC struct {
	conn	    *grpc.ClientConn
	client	    baluba.BalubaServiceClient
	chunkSize   int
	//files	    map[string]*files
	concurrency int
}

type ClientGRPCConfig struct {
	Address		string
	ChunkSize	int
	Compress	bool
	MaxConcurrency	int
}

func NewClientGRPC(cfg ClientGRPCConfig) (ClientGRPC, error) {
	var (
		client	ClientGRPC
		err	error
		grpcOpts  = []grpc.DialOption{}
	)

	if cfg.Address == "" {
		return client, errors.Errorf("Address must be specified")
	}

	if cfg.Compress {
		grpcOpts = append(grpcOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}

	grpcOpts = append(grpcOpts, grpc.WithInsecure())
	client.chunkSize = cfg.ChunkSize

	if client.conn, err = grpc.Dial(cfg.Address, grpcOpts...); err != nil {
		return client, err
	}

	client.client = baluba.NewBalubaServiceClient(client.conn)
	client.concurrency = cfg.MaxConcurrency

	return client, nil
}

func (c *ClientGRPC) Upload(ctx context.Context, dt []DirectoriesTemplate) (s Stats, err error) {
	var (
		directories []Directories
		hostname    string
	)

	if hostname, err = os.Hostname(); err != nil {
		return
	}

	if directories, err = ListFiles(dt); err != nil {
		return
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"hostname": hostname,
	}))

	if err = c.createFiles(ctx, directories); err != nil {
		return
	}

	for _, d := range directories {
		for _, f := range d.Files {
			fmt.Println(c.sendFiles(ctx, d.Path, f.Name, hostname))
		}
	}

	return
}

func (c *ClientGRPC) sendFiles(ctx context.Context, path, fname, hostname string) (err error) {
	var (
		stream	baluba.BalubaService_UploadClient
		status	*baluba.UploadStatus
		n	int
		file	*os.File
		done	= make(chan struct{})
		errc	= make(chan error, 1)
		buffer	= make([]byte, c.chunkSize)
	)

	if file, err = os.OpenFile(fmt.Sprintf("%s/%s", path, fname), os.O_RDONLY, os.ModePerm); err != nil {
		return err
	}
	defer file.Close()

	if stream, err = c.client.Upload(ctx); err != nil {
		return
	}

	go func() {
		if status, err = stream.Recv(); err != nil {
			errc <- err
			return
		}

		if status.Code != baluba.UploadStatusCode_Ok{
			errc <- errors.New(status.Message)
			return
		}

		done <- struct{}{}
		return
	}()

	go func() {
		for {
			if n, err = file.Read(buffer); err != nil {
				if err == io.EOF {
					err = nil
					break
				} else {
					errc <- err
					return
				}
			}

			if err = stream.Send(&baluba.Chunk{
				Directory:  path,
				Name:	    fname,
				Hostname:   hostname,
				Content:    buffer[:n],
			}); err != nil {
				errc <- err
				return
			}
		}
	}()

	select {
	case err = <-errc:
		return
	case _ = <-done:
		return
	}
}

func (c *ClientGRPC) createFiles(ctx context.Context, directories []Directories) (err error) {
	var (
		files	    []*baluba.Chunk
		fi	    os.FileInfo
		status	    *baluba.UploadStatus
	)

	for _, d := range directories {
		for _, f := range d.Files {
			if fi, err = os.Stat(fmt.Sprintf("%s/%s", d.Path, f.Name)); err != nil {
				continue
			}

			files = append(files, &baluba.Chunk{
				Directory:  d.Path,
				Name:	    f.Name,
				Size:	    fi.Size(),
			})
		}
	}

	if status, err = c.client.Create(ctx, &baluba.Files{
		File: files,
	}); err != nil {
		return
	}

	if status.Code == baluba.UploadStatusCode_Unknown || status.Code == baluba.UploadStatusCode_Failed {
		if status.Message == "" {
			err = errors.New("Unknow error to create files")
		} else {
			err = errors.New(status.Message)
		}
	}

	return
}

func (c *ClientGRPC) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
