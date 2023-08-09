package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type MetaphorClient struct {
	apiKey      string
	RequestBody SearchRequestBody
}

type SearchRequestBody struct {
	Query              string   `json:"query"`
	NumResults         int      `json:"numResults,omitempty"`
	IncludeDomains     []string `json:"includeDomains,omitempty"`
	ExcludeDomains     []string `json:"excludeDomains,omitempty"`
	StartCrawlDate     string   `json:"startCrawlDate,omitempty"`
	EndCrawlDate       string   `json:"endCrawlDate,omitempty"`
	StartPublishedDate string   `json:"startPublishedDate,omitempty"`
	EndPublishedDate   string   `json:"endPublishedDate,omitempty"`
	UseAutoprompt      bool     `json:"useAutoprompt,omitempty"`
	Type               string   `json:"type,omitempty"`
}

type SearchResultBody struct {
	Url           string  `json:"url"`
	Title         string  `json:"title"`
	PublishedDate string  `json:"publishedDate"`
	Author        string  `json:"author"`
	Score         float64 `json:"score"`
	Id            string  `json:"id"`
}
type SearchResponse struct {
	Results []SearchResultBody `json:"results"`
}

type ContentsResult struct {
	Url     string `json:"url"`
	Title   string `json:"title"`
	Id      string `json:"id"`
	Extract string `json:"extract"`
}

type ContentsResponse struct {
	Contents []ContentsResult `json:"contents"`
}

const (
	DefaultNumResults  = 10
	DefaultSearchURL   = "https://api.metaphor.systems/search"
	DefaultContentsURL = "https://api.metaphor.systems/contents"
)

var (
	ErrNoGoodSearchResult  = errors.New("no good search results found")
	ErrSearchRequestFailed = errors.New("search request failed")
)

func NewClient(apiKey string, options ...ClientOptions) (*MetaphorClient, error) {
	client := &MetaphorClient{
		apiKey: apiKey,
		RequestBody: SearchRequestBody{
			NumResults:     DefaultNumResults,
			IncludeDomains: []string{},
			ExcludeDomains: []string{},
			UseAutoprompt:  true,
			Type:           "neural",
		},
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}

func (client *MetaphorClient) Search(ctx context.Context, query string) (string, error) {
	client.RequestBody.Query = query

	reqBytes, err := json.Marshal(client.RequestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", DefaultSearchURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

	req.Header.Add("x-api-key", client.apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var formatedResults string
	if res.StatusCode != http.StatusOK {
		fmt.Println("status code", res.StatusCode)
		return "", ErrSearchRequestFailed
	}

	finalResults, err := client.parseResponse(ctx, body)
	if err != nil {
		return "", err
	}

	formatedResults = client.formatResults(finalResults)

	return formatedResults, nil
}

func (client *MetaphorClient) parseResponse(ctx context.Context, body []byte) ([]map[string]interface{}, error) {
	var searchResults SearchResponse
	err := json.Unmarshal(body, &searchResults)
	if err != nil {
		return nil, err
	}

	finalResults := make([]map[string]interface{}, 0)
	for _, result := range searchResults.Results {
		content := map[string]interface{}{}

		content["title"] = result.Title
		content["url"] = result.Url

		cotents, err := client.getContents(ctx, []string{result.Id})
		if err != nil {
			return finalResults, err
		}
		content["contents"] = cotents
		finalResults = append(finalResults, content)
	}

	return finalResults, nil
}

func (client *MetaphorClient) getContents(ctx context.Context, ids []string) (string, error) {
	// Convert slice of IDs to comma-separated string
	joinedIds := strings.Join(ids, "\",\"")

	// Create the dynamic URL
	url := fmt.Sprintf("%s?ids=\"%s\"", DefaultContentsURL, joinedIds)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("x-api-key", client.apiKey)
	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var contentsResults ContentsResponse
	err = json.Unmarshal(body, &contentsResults)
	if err != nil {
		return "", err
	}

	content := contentsResults.Contents[0]

	return content.Extract, nil

}

func (client *MetaphorClient) formatResults(results []map[string]interface{}) string {
	formattedResults := ""

	for _, result := range results {
		formattedResults += fmt.Sprintf("Title: %s\nContent: %s\nURL: %s\n\n", result["title"], result["contents"], result["url"])
	}

	return formattedResults
}
