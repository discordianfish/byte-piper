package pipeline

import "io"

type filter interface {
	io.Reader
	Link(r io.Reader) error
}
