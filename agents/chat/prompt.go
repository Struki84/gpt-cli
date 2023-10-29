package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/tools"
)

func Prompt(input string, options ...chains.ChainCallOption) {
	ctx := context.Background()

	llm, err := openai.NewChat(openai.WithModel("gpt-4"))
	if err != nil {
		log.Fatal(err)
	}

	// runChains(ctx, llm, input)

	runAgents(ctx, llm, input)
}

// trunk-ignore(golangci-lint/unused)
func runChains(ctx context.Context, llm llms.LanguageModel, input string) {
	chain := chains.NewConversation(llm, memory.NewConversationBuffer())

	stream := chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	})

	_, err := chains.Run(ctx, chain, input, stream)

	if err != nil {
		log.Fatal(err)
	}
}

func runAgents(ctx context.Context, llm llms.LanguageModel, input string) {

	cb := &PromptCallbacks{}

	agent := NewAsaiAgent(
		llm,
		[]tools.Tool{},
		WithCallbacksHandler(cb),
	)

	executor := NewExecutor(
		agent,
		[]tools.Tool{},
		WithCallbacksHandler(cb),
	)

	_, err := chains.Call(ctx, executor, map[string]any{"input": input})

	if err != nil {
		log.Fatal(err)
	}
}
