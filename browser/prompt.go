package browser

import (
	"context"
	"fmt"
	"gpt/util/memory"
	"gpt/util/scraper"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/duckduckgo"
	"github.com/tmc/langchaingo/tools/serpapi"
)

func Prompt(input string) string {
	var err error

	llm, err := openai.NewChat(
		openai.WithModel("gpt-3.5-turbo-16k"),
	)
	if err != nil {
		log.Print(err)
	}

	dsn := "host=localhost user=gpt-admin password=gpt-password dbname=gpt-db port=5432"
	memory := memory.NewPostgreBuffer(dsn)

	ddg, err := duckduckgo.New(5, "")
	if err != nil {
		log.Print(err)
	}

	serpapi, err := serpapi.New()
	if err != nil {
		log.Print(err)
	}

	scraper, err := scraper.NewScraper()
	if err != nil {
		log.Print(err)
	}

	tools := []tools.Tool{ddg, serpapi, scraper}

	tmpl := prompts.PromptTemplate{
		Template:       loadPromptTxToString(),
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		InputVariables: []string{"input", "history", "agent_scratchpad", "today"},
		PartialVariables: map[string]interface{}{
			"tool_names":        toolNames(tools),
			"tool_descriptions": toolDescriptions(tools),
		},
	}

	memory.SetSession("USID-002")

	ctx := context.Background()

	executor, err := agents.Initialize(
		llm,
		tools,
		agents.ZeroShotReactDescription,
		agents.WithMemory(memory),
		agents.WithPrompt(tmpl),
		agents.WithMaxIterations(3),
		// agents.WithReturnIntermediateSteps(), This throws an error(invalid input values: multiple keys and no input key set) need to open issue
	)

	if err != nil {
		log.Print(err)
	}

	answer, err := chains.Run(ctx, executor, input)
	if err != nil {
		log.Print(err)
		return ""
	}

	return answer	
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
	prompt, err := os.ReadFile("./browser/prompt.txt")
	if err != nil {
		log.Print("Error reading prompt file:", err)
	}

	return string(prompt)
}