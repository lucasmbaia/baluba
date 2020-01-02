package gogrpc

import (
	"os"
	"fmt"
	"path/filepath"
	"strings"
)

func ListFiles(dt []DirectoriesTemplate) (directories []Directories, err error) {
	checkDirectories := func(directories []Directories, path string) (int, bool) {
		for idx, d := range directories {
			if d.Path == path {
				return idx, true
			}
		}

		return 0, false
	}

	for idx, d := range dt {
		if _, err = os.Stat(d.Path); os.IsNotExist(err) || err != nil {
			return
		}

		if len(d.Files) > 0 {
			for _, file := range d.Files {
				if _, err = os.Stat(fmt.Sprintf("%s/%s", d.Path, file.Name)); os.IsNotExist(err) || err != nil {
					return
				}
			}

			directories = append(directories, Directories{
				Path:	d.Path,
				Files:	d.Files,
			})
		} else {
			if err = filepath.Walk(d.Path, func(path string, info os.FileInfo, e error) error {
				var (
					file    string
					index   int
					exists  bool
					dir     []string
				)


				if info.IsDir() {
					if _, exists = checkDirectories(directories, path); !exists {
						directories = append(directories, Directories{
							Path:   path,
						})
					}

					return nil
				}

				file = strings.Replace(path, fmt.Sprintf("%s/", d.Path), "", 1)
				if len(strings.Split(file, "/")) == 1 && !info.IsDir() {
					directories[idx].Files = append(directories[idx].Files, Files{Name: file})
				} else {
					dir = strings.Split(path, "/")

					if index, exists = checkDirectories(directories, strings.Join(dir[:len(dir) -1], "/")); exists {
						directories[index].Files = append(directories[index].Files, Files{Name: dir[len(dir)-1]})
					} else {
						directories = append(directories, Directories{
							Path:   strings.Join(dir[:len(dir) -1], "/"),
							Files:  []Files{{Name: dir[len(dir)-1]}},
						})
					}
				}

				return nil
			}); err != nil {
				return
			}
		}
	}

	return
}
