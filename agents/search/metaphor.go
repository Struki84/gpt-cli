package search

import (
	"context"
	"fmt"

	"gpt/memory"
	"gpt/tools/metaphor"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"

	options "github.com/metaphorsystems/metaphor-go"
	lc_memory "github.com/tmc/langchaingo/memory"
)

func MetaphorPrompt(input string) {
	llm, err := openai.NewChat(
		openai.WithModel("gpt-4"),
	)

	if err != nil {
		fmt.Println(err)
	}

	search, err := metaphor.NewSearch(
		options.WithAutoprompt(true),
		options.WithNumResults(5),
		options.WithType("neural"),
	)

	if err != nil {
		fmt.Println(err)
	}

	documents, err := metaphor.NewDocuments()
	if err != nil {
		fmt.Println(err)
	}

	tools := []tools.Tool{search, documents}

	dsn := "host=localhost user=gpt-admin password=gpt-password dbname=gpt-db port=5432"
	chatHistory := memory.NewPersistentChatHistory(dsn)
	agentMemory := lc_memory.NewConversationBuffer(lc_memory.WithChatHistory(chatHistory))

	chatHistory.SetSessionID("USID-003")

	tmpl := prompts.PromptTemplate{
		Template:       loadPromptTxToString(),
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		InputVariables: []string{"input", "history", "agent_scratchpad", "today"},
		PartialVariables: map[string]interface{}{
			"tool_names":        toolNames(tools),
			"tool_descriptions": toolDescriptions(tools),
		},
	}

	executor, err := agents.Initialize(
		llm,
		tools,
		agents.ZeroShotReactDescription,
		agents.WithMemory(agentMemory),
		agents.WithPrompt(tmpl),
		agents.WithMaxIterations(5),
		// agents.WithReturnIntermediateSteps(),
	)

	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()

	answer, err := chains.Run(ctx, executor, input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(answer)
}
