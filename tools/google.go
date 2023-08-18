package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const _url = "https://serpapi.com/search"

var (
	ErrNoGoodResult = errors.New("no good search results found")
	ErrAPIError     = errors.New("error from SerpAPI")
)

type Client struct {
	apiKey     string
	maxResults int
}

func NewSearch(options ...GoogleOptions) *Client {
	client := &Client{}
	for _, option := range options {
		option(client)
	}

	return client
}

func (s *Client) Search(ctx context.Context, query string) (string, error) {
	params := make(url.Values)
	query = strings.ReplaceAll(query, " ", "+")
	params.Add("q", query)
	params.Add("google_domain", "google.com")
	params.Add("gl", "us")
	params.Add("hl", "en")
	params.Add("num", fmt.Sprintf("%d", s.maxResults))
	params.Add("api_key", s.apiKey)

	reqURL := fmt.Sprintf("%s?%s", _url, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request in serpapi: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("doing response in serpapi: %w", err)
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return "", fmt.Errorf("coping data in serpapi: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		return "", fmt.Errorf("unmarshal data in serpapi: %w", err)
	}

	if errorValue, ok := result["error"]; ok {
		return "", fmt.Errorf("%w: %v", ErrAPIError, errorValue)
	}

	formattedResults := ""
	organicResults, organicResultsExists := result["organic_results"].([]interface{})

	if organicResultsExists {
		for i := 0; i < len(organicResults); i++ {
			if orgResult, ok := organicResults[i].(map[string]interface{}); ok {
				formattedResults += fmt.Sprintf("Title: %s\nDescription: %s\nURL: %s\n\n",
					orgResult["title"],
					orgResult["snippet"],
					orgResult["link"],
				)
			}
		}
	}

	return formattedResults, nil
}
