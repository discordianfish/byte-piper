package pipeline

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/docker/docker/runconfig"
)

const (
	defaultSocketPath = "/var/run/docker.sock"
	containerName     = "container.json"
)

func init() {
	inputMap["docker"] = newDockerInput
}

type dockerInput struct {
	r      *io.PipeReader
	w      *io.PipeWriter
	tw     *tar.Writer
	socket string
}

func newDockerInput(conf map[string]string) (input, error) {
	name := conf["name"]
	if name == "" {
		return nil, errors.New("name required")
	}

	r, w := io.Pipe()
	di := &dockerInput{
		r:      r,
		w:      w,
		tw:     tar.NewWriter(w),
		socket: defaultSocketPath,
	}
	if s, ok := conf["socket"]; ok {
		di.socket = s
	}

	container, containerJSON, err := di.getContainer(name)
	if err != nil {
		return nil, err
	}
	if len(container.Volumes) == 0 {
		return nil, errors.New("Container has no volumes")
	}

	go di.store(container, containerJSON) // stores files and closes writers
	return di, nil
}

func (i *dockerInput) Read(p []byte) (n int, err error) {
	return i.r.Read(p)
}

func (i *dockerInput) store(container *container, containerJSON []byte) {
	defer i.w.Close()
	defer i.tw.Close()

	now := time.Now()
	if err := i.tw.WriteHeader(&tar.Header{
		Name:       containerName,
		Size:       int64(len(containerJSON)),
		ModTime:    now,
		AccessTime: now,
		ChangeTime: now,
		Mode:       0644,
	}); err != nil {
		i.w.CloseWithError(err)
		return
	}
	if _, err := i.tw.Write(containerJSON); err != nil {
		i.w.CloseWithError(err)
		return
	}

	for _, path := range container.Volumes {
		if err := filepath.Walk(path, i.addFile); err != nil {
			i.w.CloseWithError(err)
			return
		}
	}
}

func (i *dockerInput) addFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	th, err := tar.FileInfoHeader(info, path)
	if err != nil {
		return err
	}
	th.Name = path[1:] //filepath.Join(cpath, path)
	if si, ok := info.Sys().(*syscall.Stat_t); ok {
		th.Uid = int(si.Uid)
		th.Gid = int(si.Gid)
	}

	if err := i.tw.WriteHeader(th); err != nil {
		return err
	}

	if !info.Mode().IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		if _, err := io.Copy(i.tw, file); err != nil {
			return err
		}
	}
	return nil
}

func (i *dockerInput) getContainer(name string) (*container, []byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/containers/%s/json", name), nil)
	if err != nil {
		return nil, nil, err
	}

	conn, err := net.Dial("unix", i.socket)
	if err != nil {
		return nil, nil, err
	}

	clientconn := httputil.NewClientConn(conn, nil)
	resp, err := clientconn.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		if len(body) == 0 {
			return nil, nil, fmt.Errorf("Error: %s", http.StatusText(resp.StatusCode))
		}

		return nil, nil, fmt.Errorf("HTTP %s: %s", http.StatusText(resp.StatusCode), body)
	}

	container := &container{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, body, err
	}
	return container, body, json.Unmarshal(body, &container)
}

type container struct {
	Config     runconfig.Config
	HostConfig runconfig.HostConfig
	Name       string            `json:"Name"`
	Volumes    map[string]string `json:"Volumes"`
}
