package browser

import (
	"context"
	"gpt/tools/scraper"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

func Prompt(input string) string {
	llm, err := openai.New(openai.WithModel("gpt-4"))
	if err != nil {
		return ""
	}

	ctx := context.Background()

	webDocuments, err := loadWebContent(ctx, input)
	if err != nil {
		return ""
	}

	summaryChain := chains.LoadStuffSummarization(llm)
	summary, err := chains.Call(
		ctx,
		summaryChain,
		map[string]any{"input_documents": webDocuments},
	)
	if err != nil {
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
