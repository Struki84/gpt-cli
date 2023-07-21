package chat

import (
	"context"
	"gpt/util/memory"
	"log"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/duckduckgo"
)

func Prompt(input string) string {
	llm, err := openai.New(openai.WithModel("gpt-4"))
	if err != nil {
		log.Fatal(err)
	}

	dsn := "host=localhost user=gpt-admin password=gpt-password dbname=gpt-db port=5432"
	memory := memory.NewPostgreBuffer(dsn)
	memory.SetSession("USID-001")

	ctx := context.Background()

	search, err := duckduckgo.New(5, "")
	if err != nil {
		log.Fatal(err)
	}

	agentTools := []tools.Tool{search}

	executor, err := agents.Initialize(
		llm,
		agentTools,
		agents.ConversationalReactDescription,
		agents.WithMemory(memory),
		// agents.WithReturnIntermediateSteps(), This throws an error need to open issue
	)

	if err != nil {
		log.Fatal(err)
	}

	answer, err := chains.Run(ctx, executor, input)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return answer
}
