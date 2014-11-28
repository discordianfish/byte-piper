package pipeline

import "io"

func init() {
	filterMap["rot13"] = newRot13Filter
}

type rot13Filter struct {
	r io.Reader
}

func newRot13Filter(map[string]string) (filter, error) {
	return &rot13Filter{}, nil
}

func (f *rot13Filter) Link(r io.Reader) error {
	f.r = r
	return nil
}

func (tp *rot13Filter) Read(p []byte) (n int, err error) {
	n, err = tp.r.Read(p)
	for i := 0; i < len(p); i++ {
		if (p[i] >= 'A' && p[i] < 'N') || (p[i] >= 'a' && p[i] < 'n') {
			p[i] += 13
		} else if (p[i] > 'M' && p[i] <= 'Z') || (p[i] > 'm' && p[i] <= 'z') {
			p[i] -= 13
		}
	}
	return n, err
}
