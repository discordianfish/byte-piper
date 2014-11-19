package pipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

var (
	inputMap  = make(map[string]func(conf map[string]string) (input, error))
	filterMap = make(map[string]func(conf map[string]string) (filter, error))
	outputMap = make(map[string]func(conf map[string]string) (output, error))
)

type pipeline struct {
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
func New(configFile string) (*pipeline, error) {
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
	input, err := inputNew(conf.Input.Config)
	if err != nil {
		return nil, err
	}
	outputNew, ok := outputMap[conf.Output.Type]
	if !ok {
		return nil, fmt.Errorf("Invalid output type %s", conf.Output.Type)
	}
	output, err := outputNew(conf.Output.Config)
	if err != nil {
		return nil, err
	}
	p := &pipeline{
		input:  input,
		output: output,
	}

	filterConf := &conf.Filters
	for {
		if filterConf.Type == "" {
			break
		}
		log.Printf("Filter %s", filterConf.Type)
		filterNew, ok := filterMap[filterConf.Type]
		if !ok {
			return nil, fmt.Errorf("Unknown filter %s", filterConf.Type)
		}
		filter, err := filterNew(filterConf.Config)
		if err != nil {
			return nil, err
		}
		p.filters = append(p.filters, filter)
		if filterConf.Next == nil {
			break
		}
		fc := &filterConfig{}
		if err := json.Unmarshal(filterConf.Next, fc); err != nil {
			return nil, fmt.Errorf("Coudln't unmarshal %s: %s", filterConf.Next, err)
		}
		filterConf = fc

	}
	return p, nil
}

// Run starts the pipeline.
func (p *pipeline) Run() error {
	last := p.input
	for _, f := range p.filters {
		log.Printf("Link %v -> %v", last, f)
		f.Link(last)
		last = f
	}
	_, err := io.Copy(p.output, last)
	return err
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
