package core

import (
	"io"
	"os"
	"strings"
	"path/filepath"
	"fmt"
	"sync"

	"github.com/lucasmbaia/baluba/proto"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	_ "google.golang.org/grpc/encoding/gzip"
)

type ClientGRPC struct {
	conn	  *grpc.ClientConn
	client	  baluba.BalubaServiceClient
	chunkSize int
	files	  map[string]*files
}

type ClientGRPCConfig struct {
	Address	  string
	ChunkSize int
	Compress  bool
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

	return client, nil
}

func (c *ClientGRPC) Upload(ctx context.Context, directories *[]Directories) (Stats, error) {
	var (
		stats	    Stats
		err	    error
		hostname    string
		wg	    sync.WaitGroup
		totalFiles  = 0
	)

	if err = checkDirectories(directories); err != nil {
		return stats, err
	}

	if hostname, err = os.Hostname(); err != nil {
		return stats, err
	}

	for _, d := range *directories {
		totalFiles += len(d.Files)
	}

	wg.Add(totalFiles)
	for _, d := range *directories {
		for _, f := range d.Files {
			go func(path, f string) {
				var (
					file	*os.File
					stat	os.FileInfo
					buf	= make([]byte, c.chunkSize)
		     			stream	baluba.BalubaService_UploadClient
					status	*baluba.UploadStatus
					n	int
				)

				defer wg.Done()

				if file, err = os.OpenFile(fmt.Sprintf("%s/%s", path, f), os.O_RDONLY, os.ModePerm); err != nil {
					return
				}
				defer file.Close()

				if stat, err = file.Stat(); err != nil {
					return
				}

				//if stream, err = c.client.Upload(ctx, grpc.UseCompressor("gzip")); err != nil {
				if stream, err = c.client.Upload(ctx); err != nil {
					return
				}

				if err = stream.Send(&baluba.Chunk{
					Action:	    "create_file",
					Hostname:   hostname,
					Directory:  path,
					Name:	    f,
					Size:	    stat.Size(),
				}); err != nil {
					return
				}

				if status, err = stream.Recv(); err != nil || status.Code != baluba.UploadStatusCode_Ok {
					fmt.Println(status, err)
					return
				}

				for {
					if n, err = file.Read(buf); err != nil {
						if err == io.EOF {
							break
						} else {
							return
						}
					}

					if err = stream.Send(&baluba.Chunk{
						Action:	"transfer",
						Hostname: hostname,
						Directory:  path,
						Name:	    f,
						Content:    buf[:n],
					}); err != nil {
						fmt.Println("DEU MERDA AQUI", err)
						return
					}

				}

				if status, err = stream.Recv(); err != nil || status.Code != baluba.UploadStatusCode_Ok {
					fmt.Println(status, err)
					return
				}
				/*if status, err = stream.CloseAndRecv(); err != nil {
					return
				}

				if status.Code != baluba.UploadStatusCode_Ok {
					return
				}*/

				stream.CloseSend()
				fmt.Printf("FILE  OK: %s\n", f)
				return
			}(d.Path, f.Name)
		}
	}
	wg.Wait()

	/*buf = make([]byte, c.chunkSize)
	for {
		if n, err = file.Read(buf); err != nil {
			if err == io.EOF {
				break
			} else {
				return stats, err
			}
		}

		if err = stream.Send(&baluba.Chunk{
			Content:  buf[:n],
			Name:	  "balela",
		}); err != nil {
			return stats, err
		}
	}*/

	/*if status, err = stream.CloseAndRecv(); err != nil {
		return stats, err
	}

	if status.Code != baluba.UploadStatusCode_Ok {
		return stats, errors.Errorf("upload failed - msg: %s", status.Message)
	}*/

	return stats, nil
}

func checkDirectories(directories *[]Directories) error {
	var err error

	checkDirectories := func(directories *[]Directories, path string) (int, bool) {
		for idx, d := range *directories {
			if d.Path == path {
				return idx, true
			}
		}

		return 0, false
	}

	for idx, d := range *directories {
		if _, err = os.Stat(d.Path); os.IsNotExist(err) || err != nil {
			return err
		}

		if len(d.Files) > 0 {
			for _, file := range d.Files {
				if _, err = os.Stat(fmt.Sprintf("%s/%s", d.Path, file.Name)); os.IsNotExist(err) || err != nil {
					return err
				}
			}
		} else {
			if err = filepath.Walk(d.Path, func(path string, info os.FileInfo, e error) error {
				var (
					file	string
					index	int
					exists	bool
					dir	[]string
				)


				if info.IsDir() {
					if _, exists = checkDirectories(directories, path); !exists {
						*directories = append(*directories, Directories{
							Path:	path,
						})
					}

					return nil
				}

				file = strings.Replace(path, fmt.Sprintf("%s/", d.Path), "", 1)
				if len(strings.Split(file, "/")) == 1 && !info.IsDir() {
					(*directories)[idx].Files = append((*directories)[idx].Files, Files{Name: file})
				} else {
					dir = strings.Split(path, "/")

					if index, exists = checkDirectories(directories, strings.Join(dir[:len(dir) -1], "/")); exists {
						(*directories)[index].Files = append((*directories)[index].Files, Files{Name: dir[len(dir)-1]})
					} else {
						*directories = append(*directories, Directories{
							Path:	strings.Join(dir[:len(dir) -1], "/"),
							Files:	[]Files{{Name: dir[len(dir)-1]}},
						})
					}
				}

				return nil
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ClientGRPC) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
