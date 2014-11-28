package pipeline

import (
	"compress/gzip"
	"io"
	"log"
)

func init() {
	filterMap["gzip"] = newGZipFilter
}

type gzipFilter struct {
	r io.Reader
}

func newGZipFilter(map[string]string) (filter, error) {
	return &gzipFilter{}, nil
}

func (f *gzipFilter) Link(r io.Reader) error {
	pr, pw := io.Pipe()
	f.r = pr

	zw := gzip.NewWriter(pw)
	go func() {
		defer zw.Close()
		_, err := io.Copy(zw, r)
		if err != nil {
			pr.CloseWithError(err)
			pw.CloseWithError(err)
			log.Print(err)
			return
		}
		pw.Close()
	}()
	return nil
}

func (f *gzipFilter) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}
