package pipeline

import (
	"compress/gzip"
	"io"
)

func init() {
	filterMap["gunzip"] = newGUnzipFilter
}

type gunzipFilter struct {
	r io.Reader
}

func newGUnzipFilter(map[string]string) (filter, error) {
	return &gunzipFilter{}, nil
}

func (f *gunzipFilter) Link(r io.Reader) error {
	zr, err := gzip.NewReader(r)
	f.r = zr
	return err
}

func (f *gunzipFilter) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}
