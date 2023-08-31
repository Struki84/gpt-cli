package browser

import (
	"context"
	"fmt"
	"gpt/tools/scraper"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

const SummarisationTemplate = `
	Write a concise summary of the following:

	"{{.context}}"

	CONCISE SUMMARY:`

func Prompt(input string) string {
	llm, err := openai.NewChat(openai.WithModel("gpt-4"))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	ctx := context.Background()

	webDocuments, err := loadWebContent(ctx, input)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	llmChain := chains.NewLLMChain(llm, prompts.NewPromptTemplate(
		SummarisationTemplate, []string{"context"},
	))

	summaryChain := chains.NewStuffDocuments(llmChain)
	summary, err := chains.Call(
		ctx,
		summaryChain,
		map[string]any{"input_documents": webDocuments},
	)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	response := summary["text"].(string)
	return response

}

func loadWebContent(ctx context.Context, input string) ([]schema.Document, error) {
	scraper, err := scraper.New()
	if err != nil {
		return []schema.Document{}, err
	}

	webContent, err := scraper.Call(ctx, input)
	if err != nil {
		return []schema.Document{}, err
	}

	webContentReader := strings.NewReader(webContent)

	loader := documentloaders.NewText(webContentReader)
	if err != nil {
		return []schema.Document{}, err
	}

	spliter := textsplitter.NewTokenSplitter()
	spliter.ChunkSize = 7500
	spliter.ChunkOverlap = 1024
	spliter.ModelName = "gpt-4"

	webDocuments, err := loader.LoadAndSplit(ctx, spliter)
	if err != nil {
		return []schema.Document{}, err
	}

	return webDocuments, nil
}
