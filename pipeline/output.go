package pipeline

import "io"

type output interface {
	io.Writer
}
