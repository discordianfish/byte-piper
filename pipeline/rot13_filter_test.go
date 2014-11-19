package pipeline

import (
	"bytes"
	"io"
	"testing"
)

func TestFilterRot13Chained(t *testing.T) {
	in := bytes.NewBuffer([]byte("Hello World"))
	out := &bytes.Buffer{}

	filter, err := newRot13Filter(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	filter.Link(in)
	io.Copy(out, filter)
	if out.String() != "Uryyb Jbeyq" {
		t.Fatal("Unexpected: ", out.String())
	}
}
