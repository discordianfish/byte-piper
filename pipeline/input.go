package pipeline

import "io"

type input interface {
	io.Reader
}
