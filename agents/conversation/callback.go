package conversation

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
)

const (
	DefaultFinalOutput = "AI:"
)

type PromptCallbacks struct {
	callbacks.SimpleHandler
	LastTokens             string
	PrintMode              bool
	HasKeywordBeenDetected bool
}

var _ callbacks.Handler = &PromptCallbacks{}

func (p *PromptCallbacks) HandleStreamingFunc(ctx context.Context, chunk []byte) {
	chunkStr := string(chunk)
	p.LastTokens += chunkStr

	// Buffer the last few chunks to match the keyword size
	bufferSize := len(DefaultFinalOutput)
	if len(p.LastTokens) > bufferSize {
		p.LastTokens = p.LastTokens[len(p.LastTokens)-bufferSize:]
	}

	// Check for keyword
	if p.LastTokens == DefaultFinalOutput {
		p.HasKeywordBeenDetected = true
		p.LastTokens = ""
	}

	// Check for colon and set print mode.
	if p.HasKeywordBeenDetected && chunkStr != ":" {
		p.PrintMode = true
	}

	// Print the final output after the detection of keyword.
	if p.PrintMode {
		fmt.Print(chunkStr)
	}
}
