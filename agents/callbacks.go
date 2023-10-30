package agents

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
)

var (
	DefaultKeywords = []string{"Final Answer:", "Final:", "AI:"}
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

	// Buffer the last few chunks to match the longest keyword size
	longestSize := len(DefaultKeywords[0])
	for _, k := range DefaultKeywords {
		if len(k) > longestSize {
			longestSize = len(k)
		}
	}

	if len(p.LastTokens) > longestSize {
		p.LastTokens = p.LastTokens[len(p.LastTokens)-longestSize:]
	}

	// Check for keywords
	for _, k := range DefaultKeywords {
		if p.LastTokens == k {
			p.HasKeywordBeenDetected = true
			p.LastTokens = ""
			break
		}
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
