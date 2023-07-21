package main

import (
	"context"
	"fmt"
	"io/ioutil"
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
	"gpt/util/scraper"
)

var rootCmd = &cobra.Command{
	Use:   "gpt",
	Short: "cli interface for chat gpt experimentation.",
	Long:  "Testing chat gpt.",
}

func init() {
	var config struct {
		OpenAPIKey    string `yaml:"open_api_key"`
		SerpAPIKey string `yaml:"serpapi_api_key"`
	}

	// Read the yaml file
	yamlFile, err := ioutil.ReadFile("./keys.yaml")
	if err != nil {
		panic(err)
	}

	// Unmarshal the yaml file into a Config struct
	
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

	promptCommand := &cobra.Command{
		Use: "prompt",
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]
			llm, err := openai.NewChat()
			if err != nil {
				fmt.Println(err)
			}

			completion, err := llm.Call(context.Background(), []schema.ChatMessage{
				schema.SystemChatMessage{Text: "Hello, I am a friendly chatbot. I love to talk about movies, books and music."},
				schema.HumanChatMessage{Text: input},
			}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				fmt.Print(string(chunk))
				return nil
			}))

			if err != nil {
				fmt.Println(err)
			}
	
			fmt.Println(completion)
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
			scraper, err := scraper.NewScraper()
			if err != nil {
				fmt.Println(err)
			}

			result, err := scraper.Call(context.Background(), args[0])
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(result)
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
