package chat

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type PromptCallbacks struct{}

func (*PromptCallbacks) HandleText(ctx context.Context, text string) {

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

func (*PromptCallbacks) HandleRetrieverStart(ctx context.Context, query string) {

}

func (*PromptCallbacks) HandleRetrieverEnd(ctx context.Context, query string, documents []schema.Document) {

}
