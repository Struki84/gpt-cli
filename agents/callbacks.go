package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
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

func (p *PromptCallbacks) ReadFromEgress(f func(chunk []byte)) {
	go func() {
		defer close(p.egress)
		for data := range p.egress {
			f(data)
		}
	}()
}

func (p *PromptCallbacks) HandleText(_ context.Context, text string) {
	fmt.Println(text)
}

func (p *PromptCallbacks) HandleLLMStart(_ context.Context, prompts []string) {
	fmt.Println("Entering LLM with prompts:", prompts)
}

func (p *PromptCallbacks) HandleLLMEnd(_ context.Context, output llms.LLMResult) {
	fmt.Println("Exiting LLM with results:", formatLLMResult(output))
}

func (p *PromptCallbacks) HandleChainStart(_ context.Context, inputs map[string]any) {
	fmt.Println("Entering chain with inputs:", formatChainValues(inputs))
}

func (p *PromptCallbacks) HandleChainEnd(_ context.Context, outputs map[string]any) {
	fmt.Println("Exiting chain with outputs:", formatChainValues(outputs))
}

func (p *PromptCallbacks) HandleToolStart(_ context.Context, input string) {
	fmt.Println("Entering tool with input:", removeNewLines(input))
}

func (p *PromptCallbacks) HandleToolEnd(_ context.Context, output string) {
	fmt.Println("Exiting tool with output:", removeNewLines(output))
}

func (p *PromptCallbacks) HandleAgentAction(_ context.Context, action schema.AgentAction) {
	fmt.Println("Agent selected action:", formatAgentAction(action))
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
		p.egress <- []byte(chunkStr)
	}
}

func formatChainValues(values map[string]any) string {
	output := ""
	for key, value := range values {
		output += fmt.Sprintf("\"%s\" : \"%s\", ", removeNewLines(key), removeNewLines(value))
	}

	return output
}

func formatLLMResult(output llms.LLMResult) string {
	results := "[ "
	for i := 0; i < len(output.Generations); i++ {
		for j := 0; j < len(output.Generations[i]); j++ {
			results += output.Generations[i][j].Text
		}
	}

	return results + " ]"
}

func formatAgentAction(action schema.AgentAction) string {
	return fmt.Sprintf("\"%s\" with input \"%s\"", removeNewLines(action.Tool), removeNewLines(action.ToolInput))
}

func removeNewLines(s any) string {
	return strings.ReplaceAll(fmt.Sprint(s), "\n", " ")
}
