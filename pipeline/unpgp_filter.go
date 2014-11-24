package pipeline

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"code.google.com/p/go.crypto/openpgp"
)

func init() {
	filterMap["unpgp"] = newUnpgpFilter
}

type unpgpFilter struct {
	r       io.Reader
	ready   chan bool
	keyRing openpgp.KeyRing
}

func newUnpgpFilter(conf map[string]string) (filter, error) {
	if conf["privatkey"] == "" {
		return nil, errors.New("privatkey required")
	}
	keyRing, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer([]byte(conf["privatkey"])))
	if err != nil {
		return nil, fmt.Errorf("Couldn't read privatkey: %s", err)
	}

	return &unpgpFilter{
		keyRing: keyRing,
		ready:   make(chan bool),
	}, nil
}

func (f *unpgpFilter) Link(r io.Reader) error {
	log.Print("openpgp.ReadMessage")
	message, err := openpgp.ReadMessage(r, f.keyRing, nil, nil)
	if err != nil {
		return err
	}
	f.r = message.UnverifiedBody
	return nil
}

func (f *unpgpFilter) Read(p []byte) (n int, err error) {
	if f.r == nil {
		return 0, errors.New("Couldn't decrypt message")
	}
	return f.r.Read(p)
}
