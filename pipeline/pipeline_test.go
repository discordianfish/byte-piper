package pipeline

import (
	"bytes"
	"testing"
)

func TestPipeline(t *testing.T) {
	in := bytes.NewBuffer([]byte("Hello World"))
	out := &bytes.Buffer{}

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
