package metaphor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/metaphorsystems/metaphor-go"
	"github.com/tmc/langchaingo/tools"
)

type Documents struct {
	client  *metaphor.Client
	options []metaphor.ClientOptions
}

var _ tools.Tool = &Documents{}

func NewDocuments(options ...metaphor.ClientOptions) (*Documents, error) {
	apiKey := os.Getenv("METAPHOR_API_KEY")

	client, err := metaphor.NewClient(apiKey, options...)
	if err != nil {
		return nil, err
	}

	return &Documents{
		client:  client,
		options: options,
	}, nil
}

func (tool *Documents) SetOptions(options ...metaphor.ClientOptions) {
	tool.options = options
}

func (tool *Documents) Name() string {
	return "Metaphor Contents Extractor"
}

func (tool *Documents) Description() string {
	return `
	Before using this tool use Metaphor Search and/or Metaphor Links tool.
	Metaphor Contents Extractor Retrieves contents of web pages based on a list of ID strings.
	The IDs list should be obtained from either Metaphor Search or Metaphor Search Links tool
	and are provided in the results of Metaphor Search and/or Metaphor Links tools.

	Expected input format:
	"8U71IlQ5DUTdsherhhYA,9segZCZGNjjQB2yD2uyK,..."`
}

func (tool *Documents) Call(ctx context.Context, input string) (string, error) {
	fmt.Println("Metaphor contents called with input:")
	fmt.Println("string:", input)

	ids := strings.Split(input, ",")
	for i, id := range ids {
		ids[i] = strings.TrimSpace(id)
	}

	fmt.Println("map ids:", ids)

	contents, err := tool.client.GetContents(ctx, ids)
	if err != nil {
		if errors.Is(err, metaphor.ErrNoContentExtracted) {
			return "Metaphor Extractor didn't return any results", nil
		}
		return "", err
	}

	return tool.formatContents(contents), nil
}

func (tool *Documents) formatContents(response *metaphor.ContentsResponse) string {
	formattedResults := ""

	for _, result := range response.Contents {
		formattedResults += fmt.Sprintf("Title: %s\nContent: %s\nURL: %s\n\n", result.Title, result.Extract, result.URL)
	}

	fmt.Println("Metaphor Contents Results:")
	fmt.Println(formattedResults)
	return formattedResults
}
