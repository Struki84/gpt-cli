package browser

import (
	"context"
	"fmt"
	"gpt/util/tools/metaphor"
	"gpt/util/tools/metaphor/client"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"
)

func MetaphorPrompt(input string) {
	llm, err := openai.NewChat(
		openai.WithModel("gpt-4"),
	)

	if err != nil {
		fmt.Println(err)
	}

	search, err := metaphor.NewSearch(
		client.WithAutoprompt(true),
		client.WithNumResults(5),
		client.WithType("neural"),
	)

	if err != nil {
		fmt.Println(err)
	}

	documents, err := metaphor.NewDocuments()
	if err != nil {
		fmt.Println(err)
	}

	tools := []tools.Tool{search, documents}

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
		agents.WithPrompt(tmpl),
	)

	if err != nil {
		fmt.Println(err)
	}

	answer, err := chains.Run(context.Background(), executor, input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(answer)
}
