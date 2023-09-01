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
	Please write a detailed report of the following website and its pages that will not exceed 4048 tokens:

	"{{.context}}"

	Structure the content in the following format:

	WEBSITE SUMMARY:
	[Place the summary of the entire website here]

	PAGE SUMMARIES:
	- [Page 1 Title]: [Summary of Page 1]
	- [Page N Title]: [Summary of Page N]
	...(Create a summary for every sub-page on the website)

	LINK INDEX:
	- Link 1: [Description of Link 1]
	- Link N: [Description of Link N]
	...(Depending on relevance, you can add none or N number of links)

	FINAL THOUGHTS:
	[Place any final thoughts or a concluding summary here]`

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
		chains.WithTemperature(0.1),
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
