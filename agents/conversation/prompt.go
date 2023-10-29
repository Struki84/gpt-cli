package conversation

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/tools"
)

func Prompt(input string) {
	ctx := context.Background()

	llm, err := openai.NewChat(
		openai.WithModel("gpt-4"),
	)
	if err != nil {
		fmt.Println(err)
	}

	executor, err := agents.Initialize(
		llm,
		[]tools.Tool{},
		agents.ConversationalReactDescription,
		agents.WithMemory(memory.NewSimple()),
		agents.WithCallbacksHandler(&PromptCallbacks{}),
	)

	if err != nil {
		fmt.Println(err)
	}

	_, err = chains.Run(ctx, executor, input)
	if err != nil {
		fmt.Println(err)
	}

}
