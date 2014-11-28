package pipeline

import (
	"bytes"
	"os"
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

	p := &Pipeline{
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

func TestEnvMerge(t *testing.T) {
	conf := map[string]string{
		"foo": "bar",
		"bla": "blub",
	}
	os.Setenv("FILTER_bla", "baz")
	conf = mergeEnv("FILTER_", conf)
	if conf["foo"] != "bar" || conf["bla"] != "baz" {
		t.Fatal("Unexpected conf: ", conf)
	}
}
