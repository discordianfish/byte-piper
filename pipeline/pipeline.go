package pipeline

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const defaultOutputBuffer = 1 * 1024 * 1024

var (
	inputMap  = make(map[string]func(conf map[string]string) (input, error))
	filterMap = make(map[string]func(conf map[string]string) (filter, error))
	outputMap = make(map[string]func(conf map[string]string) (output, error))

	outputBuffer = flag.Int("b", defaultOutputBuffer, "Size of output buffer")
)

type Pipeline struct {
	input   input
	filters []filter
	output  output
}

type commonConfig struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type config struct {
	Input   commonConfig `json:"input"`
	Filters filterConfig `json:"filters"`
	Output  commonConfig `json:"output"`
}

type filterConfig struct {
	commonConfig
	Next json.RawMessage `json:"next"`
}

// New returns a new pipeline.
func New(configFile string) (*Pipeline, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	conf := &config{}
	if err := json.Unmarshal(data, conf); err != nil {
		return nil, err
	}
	inputNew, ok := inputMap[conf.Input.Type]
	if !ok {
		return nil, fmt.Errorf("Invalid input type %s", conf.Input.Type)
	}
	input, err := inputNew(mergeEnv("INPUT_", conf.Input.Config))
	if err != nil {
		return nil, err
	}
	outputNew, ok := outputMap[conf.Output.Type]
	if !ok {
		return nil, fmt.Errorf("Invalid output type %s", conf.Output.Type)
	}
	output, err := outputNew(mergeEnv("OUTPUT_", conf.Output.Config))
	if err != nil {
		return nil, err
	}
	p := &Pipeline{
		input:  input,
		output: output,
	}

	filterConf := &conf.Filters
	prefix := "FILTER_"
	for {
		if filterConf.Type == "" {
			break
		}
		log.Printf("Filter %s", filterConf.Type)
		filterNew, ok := filterMap[filterConf.Type]
		if !ok {
			return nil, fmt.Errorf("Unknown filter %s", filterConf.Type)
		}
		filter, err := filterNew(mergeEnv(prefix, filterConf.Config))
		if err != nil {
			return nil, err
		}
		p.filters = append(p.filters, filter)
		if filterConf.Next == nil {
			break
		}
		fc := &filterConfig{}
		if err := json.Unmarshal(filterConf.Next, fc); err != nil {
			return nil, fmt.Errorf("Couldn't unmarshal %s: %s", filterConf.Next, err)
		}
		filterConf = fc
		prefix = prefix + "FILTER_"

	}
	return p, nil
}

// Run starts the pipeline.
func (p *Pipeline) Run() (int64, error) {
	last := p.input
	for _, f := range p.filters {
		log.Printf("Link %v -> %v", last, f)
		if err := f.Link(last); err != nil {
			return 0, err
		}
		last = f
	}
	buf := bufio.NewWriterSize(p.output, *outputBuffer)
	n, err := io.Copy(buf, last)
	if err != nil {
		return n, fmt.Errorf("Couldn't pipe data: %s", err)
	}
	log.Print("copied")
	if err := buf.Flush(); err != nil {
		return n, fmt.Errorf("Couldn't flush data: %s", err)
	}
	log.Print("flushed")
	if err := p.output.Close(); err != nil {
		return n, fmt.Errorf("Couldn't close pipeline: %s", err)
	}
	log.Print("closed")
	return n, nil
}

// Merge config with env, envs wins

// TYPE_KEY=VALUE
func mergeEnv(prefix string, conf map[string]string) map[string]string {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			panic("Invalid env variable: " + env)
		}
		if !strings.HasPrefix(parts[0], prefix) {
			continue
		}
		key := strings.TrimPrefix(parts[0], prefix)
		value := parts[1]
		conf[key] = value
	}
	return conf
}
