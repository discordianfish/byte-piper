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
	go func(tr *tar.Reader, r *io.PipeReader, w *io.PipeWriter, path string) {
		err := untar(tr, r, path)
		if err != nil {
			r.CloseWithError(err)
			w.CloseWithError(err)
		}
	}(tr, r, w, path)
	return ti, nil
}

func (o *untarOutput) Write(p []byte) (n int, err error) {
	return o.w.Write(p)
}

func (o *untarOutput) Close() error {
	return o.w.Close()
}

func untar(tr *tar.Reader, r io.Reader, path string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		path := filepath.Join(path, hdr.Name)
		info := hdr.FileInfo()
		_, errOld := os.Lstat(path)
		log.Print(path)
		if info.IsDir() {
			if os.IsNotExist(errOld) {
				if err := os.Mkdir(path, info.Mode()); err != nil {
					return err
				}
			}
			if err := os.Chown(path, hdr.Uid, hdr.Gid); err != nil {
				return err
			}
			if err := os.Chmod(path, info.Mode()); err != nil {
				return err
			}
		} else {
			file, err := os.Create(path)
			if err != nil {
				return err
			}
			if err := file.Chown(hdr.Uid, hdr.Gid); err != nil {
				return err
			}
			if err := file.Chmod(info.Mode()); err != nil {
				return err
			}
			if _, err := io.Copy(file, r); err != nil {
				return err
			}
			if err := os.Chtimes(path, hdr.AccessTime, hdr.ModTime); err != nil { // doesn't work for directories?
				return err
			}

		}
	}
}
