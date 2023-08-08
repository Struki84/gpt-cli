package metaphor

import (
	"context"
	"errors"
	"gpt/util/tools/metaphor/internal"

	"github.com/tmc/langchaingo/tools"
)

type MetaphorSearch struct {
	client *internal.MetaphorClient
}

var _ tools.Tool = &MetaphorSearch{}

func NewMetaphorSearch(apiKey string, options ...MetaphorOptions) (*MetaphorSearch, error) {
	client, err := internal.NewClient(apiKey)
	if err != nil {
		return nil, err
	}
	metaphor := &MetaphorSearch{
		client: client,
	}

	for _, option := range options {
		option(metaphor)
	}

	return metaphor, nil
}

func (tool *MetaphorSearch) Name() string {
	return "Metaphor Search"
}

func (tool *MetaphorSearch) Description() string {
	return "Metaphor Search"
}

func (tool *MetaphorSearch) Call(ctx context.Context, input string) (string, error) {
	result, err := tool.client.Search(ctx, input)
	if err != nil {
		if errors.Is(err, internal.ErrNoGoodSearchResult) {
			return "", nil
		}
		return "No good Metaphor Search Results was found", nil
	}

	return result, nil
}
