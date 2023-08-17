package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"gopkg.in/yaml.v2"

	"gpt/api"
	"gpt/browser"
	"gpt/chat"
	"gpt/read"
)

var rootCmd = &cobra.Command{
	Use:   "gpt",
	Short: "cli interface for chat gpt experimentation.",
	Long:  "Testing chat gpt.",
}

func init() {
	var config struct {
		OpenAPIKey     string `yaml:"open_api_key"`
		SerpAPIKey     string `yaml:"serpapi_api_key"`
		MetaphorAPIKey string `yaml:"metaphor_api_key"`
	}

	yamlFile, err := os.ReadFile("./keys.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	err = os.Setenv("OPENAI_API_KEY", config.OpenAPIKey)
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}

	err = os.Setenv("SERPAPI_API_KEY", config.SerpAPIKey)
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}

	err = os.Setenv("METAPHOR_API_KEY", config.MetaphorAPIKey)
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}

	promptCommand := &cobra.Command{
		Use: "prompt",
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]
			llm, err := openai.NewChat()
			if err != nil {
				fmt.Println(err)
			}

			msg := []schema.ChatMessage{
				schema.SystemChatMessage{Content: "Hello, I am a friendly chatbot. I love to talk about movies, books and music."},
				schema.HumanChatMessage{Content: input},
			}

			_, err = llm.Call(
				context.Background(),
				msg,
				llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
					fmt.Print(string(chunk))
					return nil
				}),
			)

			if err != nil {
				fmt.Println(err)
			}
		},
	}

	chatCommand := &cobra.Command{
		Use: "chat",
		Run: func(cmd *cobra.Command, args []string) {
			response := chat.Prompt(args[0])
			fmt.Println(response)
		},
	}

	searchCommand := &cobra.Command{
		Use: "browse",
		Run: func(cmd *cobra.Command, args []string) {
			response := browser.Prompt(args[0])
			fmt.Println(response)
		},
	}

	readCommand := &cobra.Command{
		Use: "read",
		Run: func(cmd *cobra.Command, args []string) {
			read.Prompt(args[0], args[1])
		},
	}

	crudCommand := &cobra.Command{
		Use: "api",
		Run: func(cmd *cobra.Command, args []string) {
			api.Prompt(args[0])
		},
	}

	runCommand := &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {

			browser.MetaphorPrompt(args[0])

			// search, err := metaphor.NewSearch(
			// 	client.WithAutoprompt(true),
			// )
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// response, err := search.Call(context.Background(), args[0])
			// if err != nil {
			// 	fmt.Print(err)
			// 	return
			// }

			// fmt.Println("Search Results:")
			// fmt.Println(response)

			// document, err := metaphor.NewDocuments()
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// contents, err := document.Call(context.Background(), "NlmibOCRhQqIaknRQqvWtQ")
			// if err != nil {
			// 	fmt.Print(err)
			// 	return
			// }

			// fmt.Println("RDJ Twitter contents:")
			// fmt.Println(contents)

			// searchLinks, err := metaphor.NewLinksSearch(
			// 	client.WithNumResults(5),
			// )

			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// similarLinks, err := searchLinks.Call(context.Background(), "https://www.wikidata.org/wiki/Q587102")
			// if err != nil {
			// 	fmt.Print(err)
			// 	return
			// }

			// fmt.Println("Similar links:")
			// fmt.Println(similarLinks)

			// key := os.Getenv("SERPAPI_API_KEY")
			// search := tools.NewSearch(
			// 	tools.WithApiKey(key),
			// 	tools.WithMaxResults(5),
			// )

			// result, err := search.Search(context.Background(), args[0])
			// if err != nil {
			// 	fmt.Print(err)
			// }

			// scraper, err := scraper.New(
			// 	scraper.WithAsync(false),
			// )
			// if err != nil {
			// 	fmt.Println(err)
			// }

			// result, err := scraper.Call(context.Background(), args[0])
			// if err != nil {
			// 	fmt.Println(err)
			// }

			// fmt.Println(result)
		},
	}

	rootCmd.AddCommand(promptCommand)
	rootCmd.AddCommand(chatCommand)
	rootCmd.AddCommand(searchCommand)
	rootCmd.AddCommand(readCommand)
	rootCmd.AddCommand(crudCommand)
	rootCmd.AddCommand(runCommand)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
