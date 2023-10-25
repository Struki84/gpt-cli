package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
)

func Prompt(input string, options ...chains.ChainCallOption) {
	ctx := context.Background()

	llm, err := openai.NewChat(openai.WithModel("gpt-4"))
	if err != nil {
		log.Fatal(err)
	}

	chain := chains.NewConversation(llm, memory.NewConversationBuffer())

	stream := chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Println(string(chunk))
		return nil
	})

	_, err = chains.Run(ctx, chain, input, stream)

	if err != nil {
		log.Fatal(err)
	}

	// executor, err := agents.Initialize(
	// 	llm,
	// 	[]tools.Tool{},
	// 	agents.ConversationalReactDescription,
	// 	agents.WithReturnIntermediateSteps(), // This throws an error
	// 	agents.WithCallbacksHandler(&PromptCallbacks{}),
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = chains.Run(
	// 	ctx,
	// 	executor,
	// 	input,
	// 	chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	// 		fmt.Println(string(chunk))
	// 		return nil
	// 	}),
	// )

	if err != nil {
		log.Fatal(err)
	}
}
