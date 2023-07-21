package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
)

func Prompt(query string) {
	
	llm, err := openai.New()
	if err != nil {
		log.Print(err)
	}

	apiDocs := loadApiDocs()
	chain := NewAPIChain(llm, http.DefaultClient)
	
	input := map[string]any {
		"api_docs": apiDocs,
		"input":  query,
	}
	
	result, err := chains.Call(context.Background(), chain, input, chains.WithTemperature(0.1))
	if err != nil {
		log.Print(err)
	}
	
	fmt.Println(result["answer"])
}

func loadApiDocs() string {
	prompt, err := os.ReadFile("./api/docs/mockup_api_docs.txt")
	if err != nil {
		log.Print("Error reading api docs:", err)
		return ""
	}

	return string(prompt)
}