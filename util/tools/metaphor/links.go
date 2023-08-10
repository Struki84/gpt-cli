package metaphor

import (
	"context"
	"errors"
	"fmt"
	"gpt/util/tools/metaphor/internal"
	"os"

	"github.com/tmc/langchaingo/tools"
)

type MetaphorLinksSearch struct {
	client *internal.MetaphorClient
}

var _ tools.Tool = &MetaphorLinksSearch{}

func NewLinksSearch(options ...internal.ClientOptions) (*MetaphorLinksSearch, error) {
	apiKey := os.Getenv("METAPHOR_API_KEY")
	if apiKey == "" {
		return nil, ErrMissingToken
	}

	client, err := internal.NewClient(apiKey, options...)
	if err != nil {
		return nil, err
	}
	metaphor := &MetaphorLinksSearch{
		client: client,
	}

	return metaphor, nil
}

func (tool *MetaphorLinksSearch) Name() string {
	return "Metaphor Links Search"
}

func (tool *MetaphorLinksSearch) Description() string {
	return `
	Metaphor Links Search finds similar links to the link provided.
	Input should be the url for which you would like to find similar links`
}

func (tool *MetaphorLinksSearch) Call(ctx context.Context, input string) (string, error) {
	links, err := tool.client.FindSimilar(ctx, input)
	if err != nil {
		if errors.Is(err, internal.ErrNoLinksFound) {
			return "Metaphor Links Search didn't return any results", nil
		}
		return "", err
	}

	return tool.formatLinks(links), nil
}

func (tool *MetaphorLinksSearch) formatLinks(response *internal.SearchResponse) string {
	formattedResults := ""

	for _, result := range response.Results {
		formattedResults += fmt.Sprintf("Title: %s\nURL: %s\n\n", result.Title, result.Url)
	}

	return formattedResults
}
