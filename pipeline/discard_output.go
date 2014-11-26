package pipeline

func init() {
	outputMap["discard"] = newDiscardOutput
}

type discardOutput struct{}

func newDiscardOutput(conf map[string]string) (output, error) {
	return &discardOutput{}, nil
}

func (o *discardOutput) Write(p []byte) (int, error) {
	return len(p), nil
}

func (o *discardOutput) Close() error { return nil }
