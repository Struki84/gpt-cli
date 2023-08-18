package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
)

const (
	// nolint: lll
	_llmAPIURLPrompt = `
	You are given the API Documentation:

	{{.api_docs}}

	Your task is to construct a full API JSON object based on the provided input. The input could be a question that requires an API call for its answer, or a direct or indirect instruction to consume an API. The input will be unpredictable and could come from a user or an agent.

	Your goal is to create an API call that accurately reflects the intent of the input. Be sure to exclude any unnecessary data in the API call to ensure efficiency.

	Input: {{.input}}

	Respond with a JSON object.

	{
		"method":  [the HTTP method for the API call, such as GET or POST],
		"headers": [object representing the HTTP headers required for the API call, always add a "Content-Type" header],
		"url": 	   [full for the API call],
		"body":    [object containing the data sent with the request, if needed]
	}`

	// nolint: lll
	_llmAPIResponsePrompt = _llmAPIURLPrompt + `
	Here is the response from the API:

	{{.api_response}}

	Now, summarize this response. Your summary should reflect the original input and highlight the key information from the API response that answers or relates to that input. Try to make your summary concise, yet informative.

	Summary:`
)

type HTTPRequest interface {
	Do(*http.Request) (*http.Response, error)
}

type APIChain struct {
	RequestChain *chains.LLMChain
	AnswerChain  *chains.LLMChain
	Request    HTTPRequest
}

func NewAPIChain(llm llms.LanguageModel, request HTTPRequest) APIChain {
	reqPrompt := prompts.NewPromptTemplate(_llmAPIURLPrompt, []string{"api_docs", "input"})
	reqChain := chains.NewLLMChain(llm, reqPrompt)

	resPrompt := prompts.NewPromptTemplate(_llmAPIResponsePrompt, []string{"input", "api_docs", "api_response"})
	resChain := chains.NewLLMChain(llm, resPrompt)

	return APIChain{
		RequestChain: reqChain,
		AnswerChain:  resChain,
		Request:      request,
	}
}

func (chain APIChain) Call(ctx context.Context, values map[string]any, opts ...chains.ChainCallOption) (map[string]any, error) {
	fmt.Println("Request Chain Calling...")

	var output struct {
		Method  string            `json:"method"`
		Headers map[string]string `json:"headers"`
		URL     string            `json:"url"`
		Body    map[string]string `json:"body"`
	}

	opts = append(opts, chains.WithTemperature(0.1))
	
	tmpOutput, err := chains.Call(ctx, chain.RequestChain, values, opts...)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(?s)\{.*\}`)
	jsonString := re.FindString(tmpOutput["text"].(string))

	fmt.Println("Request Chain Output JSON string: ")
	fmt.Println(jsonString)

	// Convert the LLM output into the anonymous struct.
	err = json.Unmarshal([]byte(jsonString), &output)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	apiResponse, err := chain.runRequest(output.Method, output.URL, output.Headers, output.Body)
	if err != nil {
		return nil, err
	}

	tmpOutput["input"] = values["input"]
	tmpOutput["api_docs"] = values["api_docs"]
	tmpOutput["api_response"] = apiResponse

	fmt.Println("Api Response: ")
	fmt.Println(apiResponse)

	answer, err := chains.Call(ctx, chain.AnswerChain, tmpOutput, opts...)
	if err != nil {
		return nil, err
	}

	return map[string]any{"answer": answer["text"]}, err
}

func (a APIChain) GetMemory() schema.Memory { //nolint:ireturn
	return memory.NewSimple()
}

func (a APIChain) GetInputKeys() []string {
	return []string{"api_docs", "input"}
}

func (a APIChain) GetOutputKeys() []string {
	return []string{"answer"}
}

func (chain APIChain) runRequest(method string, url string, headers map[string]string, body map[string]string) (string, error) {
	var bodyReader io.Reader
	
	if method == "POST" || method == "PUT" {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return "", err
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	// Create the new request defined by reqChain
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
        return "", err
    }

	// set request headers passed from reqChain
	for key, value := range headers {
		req.Header.Add(key, value)
	}
    
	resp, err := chain.Request.Do(req)
	if err != nil {
		return "", err
	}
	
	defer resp.Body.Close()
	
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(b), nil
}