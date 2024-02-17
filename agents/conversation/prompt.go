package conversation

import (
	"context"
	"fmt"

	my_agents "gpt/agents"
	"gpt/tools/documents"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/tools"
)

func Prompt(input string) {
	agentCallback := my_agents.NewPromptCallbacks()

	agentCallback.ReadFromEgress(func(chunk []byte) {
		fmt.Print(string(chunk))
	})

	ctx := context.Background()

	llm, err := openai.NewChat(
		openai.WithModel("gpt-4"),
	)

	if err != nil {
		fmt.Println(err)
	}

	docTool, err := documents.NewLoader()
	if err != nil {
		fmt.Println(err)
		return
	}

	executor, err := agents.Initialize(
		llm,
		[]tools.Tool{docTool},
		agents.ConversationalReactDescription,
		agents.WithMemory(memory.NewSimple()),
		agents.WithCallbacksHandler(callbacks.LogHandler{}),
	)

	if err != nil {
		fmt.Println(err)
	}

	response, err := chains.Run(ctx, executor, input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response)
}
