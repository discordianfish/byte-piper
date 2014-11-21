package pipeline

import (
	"bytes"
	"testing"
)

// dummy buffer for testing
type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() (err error) {
	return
}

func TestPipeline(t *testing.T) {
	in := bytes.NewBuffer([]byte("Hello World"))
	out := &buffer{} //bytes.Buffer{}

	p := &pipeline{
		input:  in,
		output: out,
	}
	if err := p.Run(); err != nil {
		t.Fatal(err)
	}
	if out.String() != "Hello World" {
		t.Fatal("Unexpected: ", out.String())
	}
}
