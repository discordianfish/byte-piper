package pipeline

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

type tarInput struct {
	path      string
	r         io.Reader
	tarWriter *tar.Writer
}

func init() {
	inputMap["tar"] = newTarInput
}

func newTarInput(conf map[string]string) (input, error) {
	path := conf["path"]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist", path)
	}

	r, w := io.Pipe()
	tarWriter := tar.NewWriter(w)
	ti := &tarInput{
		path:      path,
		tarWriter: tarWriter,
		r:         r,
	}
	go func(w io.WriteCloser) {
		defer ti.tarWriter.Close() // This doesn *not* close the embedded writer
		defer w.Close()            // So we do it here
		if err := filepath.Walk(path, ti.addFile); err != nil {
			log.Printf("Couldn't walk %s: %s", path, err)
		}
	}(w)
	return ti, nil
}

func (i *tarInput) Read(p []byte) (n int, err error) {
	return i.r.Read(p)
}

func (i *tarInput) addFile(path string, info os.FileInfo, err error) error {
	log.Printf("file %s", path)
	if err != nil {
		return err
	}

	relPath := path[len(filepath.Dir(i.path))+1:] // relative to volume parent directory (<docker>/vfs/dir/)
	th, err := tar.FileInfoHeader(info, relPath)
	if err != nil {
		return err
	}
	th.Name = relPath
	if si, ok := info.Sys().(*syscall.Stat_t); ok {
		th.Uid = int(si.Uid)
		th.Gid = int(si.Gid)
	}

	if err := i.tarWriter.WriteHeader(th); err != nil {
		return err
	}

	if !info.Mode().IsDir() && th.Typeflag != tar.TypeSymlink {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		log.Printf("copying file %s", path)
		if _, err := io.Copy(i.tarWriter, file); err != nil {
			return err
		}

		log.Printf("done!")
	}
	return nil

}
