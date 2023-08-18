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

	llm, err := openai.New(openai.WithModel("gpt-4"))
	if err != nil {
		log.Print(err)
	}

	apiDocs := loadApiDocs("mockup_api_docs.txt")
	chain := NewAPIChain(llm, http.DefaultClient)

	input := map[string]any{
		"api_docs": apiDocs,
		"input":    query,
	}

	result, err := chains.Call(context.Background(), chain, input)
	if err != nil {
		log.Print(err)
	}

	fmt.Println(result["answer"])
}

func loadApiDocs(file string) string {
	path := "./api/docs/" + file
	prompt, err := os.ReadFile(path)
	if err != nil {
		log.Print("Error reading api docs:", err)
		return ""
	}

	return string(prompt)
}
