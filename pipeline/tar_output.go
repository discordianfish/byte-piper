package pipeline

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type untarOutput struct {
	path string
	tr   *tar.Reader
	w    io.WriteCloser
}

func init() {
	outputMap["untar"] = newUntarOutput
}

func newUntarOutput(conf map[string]string) (output, error) {
	path := conf["path"]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist", path)
	}

	r, w := io.Pipe()
	tr := tar.NewReader(r)
	ti := &untarOutput{
		path: path,
		w:    w,
		tr:   tr,
	}
	go func() {
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Print(err)
				return
			}
			path := filepath.Join(path, hdr.Name)
			info := hdr.FileInfo()
			if info.IsDir() {
				if err := os.Mkdir(path, info.Mode()); err != nil {
					log.Print(err)
					return
				}
				if err := os.Chown(path, hdr.Uid, hdr.Gid); err != nil {
					log.Print(err)
					return
				}
				if err := os.Chmod(path, info.Mode()); err != nil {
					log.Print(err)
					return
				}
			} else {
				file, err := os.Create(path)
				if err != nil {
					log.Print(err)
					return
				}
				if err := file.Chown(hdr.Uid, hdr.Gid); err != nil {
					log.Print(err)
					return
				}
				if err := file.Chmod(info.Mode()); err != nil {
					log.Print(err)
					return
				}
				if _, err := io.Copy(file, r); err != nil {
					log.Print(err)
					return
				}
				if err := os.Chtimes(path, hdr.AccessTime, hdr.ModTime); err != nil { // doesn't work for directories?
					log.Print(err)
					return
				}

			}
		}
	}()
	return ti, nil
}

func (o *untarOutput) Write(p []byte) (n int, err error) {
	return o.w.Write(p)
}

func (o *untarOutput) Close() error {
	return o.w.Close()
}
