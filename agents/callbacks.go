package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
)

var (
	DefaultKeywords = []string{"Final Answer:", "Final:", "AI:"}
)

type PromptCallbacks struct {
	callbacks.SimpleHandler
	egress          chan []byte
	Keywords        []string
	LastTokens      string
	KeywordDetected bool
	PrintOutput     bool
}

var _ callbacks.Handler = &PromptCallbacks{}

func NewPromptCallbacks(keywords ...string) *PromptCallbacks {
	if len(keywords) > 0 {
		DefaultKeywords = keywords
	}

	return &PromptCallbacks{
		egress:   make(chan []byte),
		Keywords: DefaultKeywords,
	}
}

func (p *PromptCallbacks) GetEgress() chan []byte {
	return p.egress
}

func (p *PromptCallbacks) HandleChainStart(ctx context.Context, inputs map[string]any) {
	fmt.Println("Chain Starting...")
}

func (p *PromptCallbacks) HandleStreamingFunc(ctx context.Context, chunk []byte) {
	chunkStr := string(chunk)
	p.LastTokens += chunkStr

	// Buffer the last few chunks to match the longest keyword size
	longestSize := len(p.Keywords[0])
	for _, k := range p.Keywords {
		if len(k) > longestSize {
			longestSize = len(k)
		}
	}

	if len(p.LastTokens) > longestSize {
		p.LastTokens = p.LastTokens[len(p.LastTokens)-longestSize:]
	}

	// Check for keywords
	for _, k := range DefaultKeywords {
		if strings.Contains(p.LastTokens, k) {
			p.KeywordDetected = true
		}
	}

	// Check for colon and set print mode.
	if p.KeywordDetected && chunkStr != ":" {
		p.PrintOutput = true
	}

	// Print the final output after the detection of keyword.
	if p.PrintOutput {
		p.egress <- chunk
	}
}
