package chat

import (
	"context"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

var _ callbacks.Handler = &PromptCallbacks{}

type PromptCallbacks struct {
	callbacks.SimpleHandler
}

func (*PromptCallbacks) HandleStreamingFunc(ctx context.Context, chunk []byte) {
	print(string(chunk))
}

func (*PromptCallbacks) HandleLLMStart(ctx context.Context, prompts []string) {

}

func (*PromptCallbacks) HandleLLMEnd(ctx context.Context, output llms.LLMResult) {

}

func (*PromptCallbacks) HandleChainStart(ctx context.Context, inputs map[string]any) {
	println("HandleChainStart")
}

func (*PromptCallbacks) HandleChainEnd(ctx context.Context, outputs map[string]any) {

}

func (*PromptCallbacks) HandleToolStart(ctx context.Context, input string) {

}

func (*PromptCallbacks) HandleToolEnd(ctx context.Context, output string) {

}

func (*PromptCallbacks) HandleAgentAction(ctx context.Context, action schema.AgentAction) {
	println("Agent is performing an action")
}
