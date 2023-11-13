package ollama

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func Prompt(input string) {

	llm, err := ollama.New(
		ollama.WithModel("mistral"),
		ollama.WithServerURL("http://localhost:11434"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	completion, err := llm.Call(ctx,
		input,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = completion

}
