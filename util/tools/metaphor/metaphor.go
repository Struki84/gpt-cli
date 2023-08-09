package metaphor

import (
	"context"
	"errors"
	"fmt"
	"gpt/util/tools/metaphor/internal"
	"os"

	"github.com/tmc/langchaingo/tools"
)

type MetaphorSearch struct {
	client *internal.MetaphorClient
}

var (
	_ tools.Tool = &MetaphorSearch{}

	ErrMissingToken = errors.New("missing the Metaphor API key, set it in the METAPHOR_API_KEY environment variable")
)

func NewSearch(options ...internal.ClientOptions) (*MetaphorSearch, error) {
	apiKey := os.Getenv("METAPHOR_API_KEY")
	if apiKey == "" {
		return nil, ErrMissingToken
	}

	client, err := internal.NewClient(apiKey, options...)
	if err != nil {
		return nil, err
	}
	metaphor := &MetaphorSearch{
		client: client,
	}

	return metaphor, nil
}

func (tool *MetaphorSearch) Name() string {
	return "Metaphor Search"
}

func (tool *MetaphorSearch) Description() string {
	return `
	Metaphor Search uses a transformer architecture to predict links given text,
	and it gets its power from having been trained on the way that people talk
	about links on the Internet. This training produces a model that returns
	links that are both high in relevance and quality. However, the model does
	expect queries that look like how people describe a link on the Internet.
	For example:
	"'best restaurants in SF" is a bad query, whereas
	"Here is the best restaurant in SF:" is a good query.
	`
}

func (tool *MetaphorSearch) Call(ctx context.Context, input string) (string, error) {
	result, err := tool.client.Search(ctx, input)
	if err != nil {
		if errors.Is(err, internal.ErrNoGoodSearchResult) {
			return "No good Metaphor Search Results was found", nil
		}
		fmt.Println("Call error", err)
		return "", nil
	}

	return result, nil
}
