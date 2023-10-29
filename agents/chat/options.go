package chat

import (
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

type CreationOptions struct {
	prompt                  prompts.PromptTemplate
	memory                  schema.Memory
	callbacksHandler        callbacks.Handler
	maxIterations           int
	returnIntermediateSteps bool
	outputKey               string
	promptPrefix            string
	formatInstructions      string
	promptSuffix            string
}

// CreationOption is a function type that can be used to modify the creation of the agents
// and executors.
type CreationOption func(*CreationOptions)

func executorDefaultOptions() CreationOptions {
	return CreationOptions{
		maxIterations: 5,
		outputKey:     "output",
		memory:        memory.NewSimple(),
	}
}

func conversationalDefaultOptions() CreationOptions {
	return CreationOptions{
		promptPrefix:       _defaultConversationalPrefix,
		formatInstructions: _defaultConversationalFormatInstructions,
		promptSuffix:       _defaultConversationalSuffix,
		outputKey:          "output",
	}
}

func (co CreationOptions) getConversationalPrompt(tools []tools.Tool) prompts.PromptTemplate {
	if co.prompt.Template != "" {
		return co.prompt
	}

	return createConversationalPrompt(
		tools,
		co.promptPrefix,
		co.formatInstructions,
		co.promptSuffix,
	)
}

func WithMaxIterations(iterations int) CreationOption {
	return func(co *CreationOptions) {
		co.maxIterations = iterations
	}
}

// WithOutputKey is an option for setting the output key of the agent.
func WithOutputKey(outputKey string) CreationOption {
	return func(co *CreationOptions) {
		co.outputKey = outputKey
	}
}

// WithReturnIntermediateSteps is an option for making the executor return the intermediate steps
// taken.
func WithReturnIntermediateSteps() CreationOption {
	return func(co *CreationOptions) {
		co.returnIntermediateSteps = true
	}
}

func WithMemory(m schema.Memory) CreationOption {
	return func(co *CreationOptions) {
		co.memory = m
	}
}

func WithCallbacksHandler(handler callbacks.Handler) CreationOption {
	return func(co *CreationOptions) {
		co.callbacksHandler = handler
	}
}
