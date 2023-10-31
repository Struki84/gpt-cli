package chat

import (
	"context"
	"fmt"
	"log"

	my_agents "gpt/agents"

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

// nolint: all
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

	agentCallback := my_agents.NewPromptCallbacks()

	agent := NewAsaiAgent(
		llm,
		[]tools.Tool{},
		WithCallbacksHandler(agentCallback),
	)

	executor := NewExecutor(
		agent,
		[]tools.Tool{},
		WithCallbacksHandler(agentCallback),
	)

	_, err := chains.Call(ctx, executor, map[string]any{"input": input})

	if err != nil {
		log.Fatal(err)
	}

	egressChannel := agentCallback.GetEgress()

	go func() {
		for data := range egressChannel {
			fmt.Print(string(data))
		}
	}()
}
