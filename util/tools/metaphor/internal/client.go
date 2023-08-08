package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type MetaphorClient struct {
	apiKey      string
	RequestBody SearchRequestBody
}

type SearchRequestBody struct {
	Query              string   `json:"query"`
	NumResults         int      `json:"numResults"`
	IncludeDomains     []string `json:"includeDomains"`
	ExcludeDomains     []string `json:"excludeDomains"`
	StartCrawlDate     string   `json:"startCrawlDate"`
	EndCrawlDate       string   `json:"endCrawlDate"`
	StartPublishedDate string   `json:"startPublishedDate"`
	EndPublishedDate   string   `json:"endPublishedDate"`
	UseAutoprompt      bool     `json:"useAutoprompt"`
	Type               string   `json:"type"`
}

const (
	DefualtNumResults = 10
	DefaultSearchURL  = "https://api.metaphor.systems/search"
)

var ErrNoGoodSearchResult = errors.New("no good search results found")

func NewClient(apiKey string) (*MetaphorClient, error) {
	client := &MetaphorClient{
		apiKey: apiKey,
		RequestBody: SearchRequestBody{
			Query:              "",
			NumResults:         DefualtNumResults,
			IncludeDomains:     []string{},
			ExcludeDomains:     []string{},
			StartCrawlDate:     "",
			EndCrawlDate:       "",
			StartPublishedDate: "",
			EndPublishedDate:   "",
			UseAutoprompt:      false,
			Type:               "",
		},
	}

	return client, nil
}

func (c *MetaphorClient) Search(ctx context.Context, query string) (string, error) {

	reqBytes, err := json.Marshal(c.RequestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", DefaultSearchURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

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

	response := string(body)

	fmt.Println(response)
	return response, nil
}
