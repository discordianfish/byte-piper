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
	filterMap["pgp"] = newPGPFilter
}

type pgpFilter struct {
	r     io.Reader
	pw    io.WriteCloser
	pgpw  io.WriteCloser
	ready chan bool
}

func newPGPFilter(conf map[string]string) (filter, error) {
	if conf["pubkey"] == "" {
		return nil, errors.New("pubkey required")
	}
	to, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer([]byte(conf["pubkey"])))
	if err != nil {
		return nil, fmt.Errorf("Couldn't read pubkey: %s", err)
	}

	pr, pw := io.Pipe()
	f := &pgpFilter{
		r:     pr,
		pw:    pw,
		ready: make(chan bool),
	}
	go func() {
		w, err := openpgp.Encrypt(pw, to, nil, nil, nil) // Not signing yet, sorry
		f.pgpw = w
		f.ready <- true
		if err != nil {
			log.Print(err)
			return
		}
	}()
	return f, nil
}

func (f *pgpFilter) Link(r io.Reader) {
	go func() {
		// openpgp.Encrypt() blocks and waits for some input,
		// so we need to wait for it before we can read from
		// the writher
		<-f.ready
		defer f.pw.Close()
		defer f.pgpw.Close()
		if _, err := io.Copy(f.pgpw, r); err != nil {
			log.Print(err)
			return
		}
	}()
}

func (f *pgpFilter) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}
