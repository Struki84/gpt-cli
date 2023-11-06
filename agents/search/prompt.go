package search

import (
	"context"
	"fmt"
	"gpt/tools/scraper"
	"log"
	"os"
	"strings"

	my_agents "gpt/agents"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	lc_memory "github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/duckduckgo"
	"github.com/tmc/langchaingo/tools/serpapi"
)

func Prompt(input string) {
	var err error

	llm, err := openai.NewChat(
		openai.WithModel("gpt-4"),
	)
	if err != nil {
		log.Print(err)
	}

	// dsn := "host=localhost user=gpt-admin password=gpt-password dbname=gpt-db port=5432"

	// chatHistory := memory.NewPersistentChatHistory(dsn)
	// chatHistory.SetSessionID("USID-003")

	// agentMemory := lc_memory.NewConversationBuffer(lc_memory.WithChatHistory(chatHistory))

	agentMemory := lc_memory.NewSimple()

	ddg, err := duckduckgo.New(5, duckduckgo.DefaultUserAgent)
	if err != nil {
		log.Print(err)
	}

	serpapi, err := serpapi.New()
	if err != nil {
		log.Print(err)
	}

	scraper, err := scraper.New()
	if err != nil {
		log.Print(err)
	}

	tools := []tools.Tool{ddg, serpapi, scraper}

	tmpl := prompts.PromptTemplate{
		Template:       loadPromptTxToString(),
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		InputVariables: []string{"input", "agent_scratchpad", "today"},
		PartialVariables: map[string]interface{}{
			"tool_names":        toolNames(tools),
			"tool_descriptions": toolDescriptions(tools),
			"history":           "",
		},
	}

	ctx := context.Background()

	agentCallback := my_agents.NewPromptCallbacks()

	agentCallback.ReadFromEgress(func(chunk []byte) {
		fmt.Print(string(chunk))
	})

	executor, err := agents.Initialize(
		llm,
		tools,
		agents.ZeroShotReactDescription,
		agents.WithMemory(agentMemory),
		agents.WithPrompt(tmpl),
		agents.WithMaxIterations(3),
		agents.WithCallbacksHandler(agentCallback),
	)

	if err != nil {
		log.Print(err)
	}

	inputs := map[string]any{
		"input": input,
	}

	_, err = chains.Predict(ctx, executor, inputs)
	if err != nil {
		log.Print(err)
	}
}

func toolNames(tools []tools.Tool) string {
	var tn strings.Builder
	for i, tool := range tools {
		if i > 0 {
			tn.WriteString(", ")
		}
		tn.WriteString(tool.Name())
	}

	return tn.String()
}

func toolDescriptions(tools []tools.Tool) string {
	var ts strings.Builder
	for _, tool := range tools {
		ts.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name(), tool.Description()))
	}

	return ts.String()
}

func loadPromptTxToString() string {
	prompt, err := os.ReadFile("./agents/search/prompt.txt")
	if err != nil {
		log.Print("Error reading prompt file:", err)
	}

	return string(prompt)
}
