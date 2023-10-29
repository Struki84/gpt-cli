package chat

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

const (
	ConversationFinalAnswerAction = "AI"
)

var (
	// ErrExecutorInputNotString is returned if an input to the executor call function is not a string.
	ErrExecutorInputNotString = errors.New("input to executor not string")
	// ErrAgentNoReturn is returned if the agent returns no actions and no finish.
	ErrAgentNoReturn = errors.New("no actions or finish was returned by the agent")
	// ErrNotFinished is returned if the agent does not give a finish before  the number of iterations
	// is larger then max iterations.
	ErrNotFinished = errors.New("agent not finished before max iterations")
	// ErrUnknownAgentType is returned if the type given to the initializer is invalid.
	ErrUnknownAgentType = errors.New("unknown agent type")
	// ErrInvalidOptions is returned if the options given to the initializer is invalid.
	ErrInvalidOptions = errors.New("invalid options")

	// ErrUnableToParseOutput is returned if the output of the llm is unparsable.
	ErrUnableToParseOutput = errors.New("unable to parse agent output")
	// ErrInvalidChainReturnType is returned if the internal chain of the agent eturns a value in the
	// "text" filed that is not a string.
	ErrInvalidChainReturnType = errors.New("agent chain did not return a string")
)

type StreamingFunc func(ctx context.Context, chunk []byte) error

var _ agents.Agent = &AsaiAgent{}

type AsaiAgent struct {
	Chain            chains.Chain
	Tools            []tools.Tool
	OutputKey        string
	CallbacksHandler callbacks.Handler
}

func NewAsaiAgent(llm llms.LanguageModel, tools []tools.Tool, opts ...CreationOption) *AsaiAgent {
	options := conversationalDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &AsaiAgent{
		Chain:            chains.NewLLMChain(llm, options.getConversationalPrompt(tools)),
		Tools:            tools,
		OutputKey:        "output",
		CallbacksHandler: options.callbacksHandler,
	}
}

func (agent *AsaiAgent) Plan(ctx context.Context, intermediateSteps []schema.AgentStep, inputs map[string]string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	fullInputs := make(map[string]any, len(inputs))
	for k, v := range inputs {
		fullInputs[k] = v
	}

	fullInputs["agent_scratchpad"] = constructScratchPad(intermediateSteps)

	var stream func(ctx context.Context, chunk []byte) error

	if agent.CallbacksHandler != nil {
		stream = func(ctx context.Context, chunk []byte) error {
			agent.CallbacksHandler.HandleStreamingFunc(ctx, chunk)
			return nil
		}
	}

	output, err := chains.Predict(
		ctx,
		agent.Chain,
		fullInputs,
		chains.WithStopWords([]string{"\nObservation:", "\n\tObservation:"}),
		chains.WithStreamingFunc(stream),
	)

	if err != nil {
		return nil, nil, err
	}

	return agent.parseOutput(output)
}

func (agent *AsaiAgent) GetInputKeys() []string {
	return []string{}
}

func (agent *AsaiAgent) GetOutputKeys() []string {
	return []string{agent.OutputKey}
}

func (agent *AsaiAgent) parseOutput(output string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	if strings.Contains(output, ConversationFinalAnswerAction) {
		splits := strings.Split(output, ConversationFinalAnswerAction)

		return nil, &schema.AgentFinish{
			ReturnValues: map[string]any{
				agent.OutputKey: splits[len(splits)-1],
			},
			Log: output,
		}, nil
	}

	r := regexp.MustCompile(`Action: (.*?)[\n]*Action Input: (.*)`)
	matches := r.FindStringSubmatch(output)
	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("%w: %s", ErrUnableToParseOutput, output)
	}

	return []schema.AgentAction{
		{Tool: strings.TrimSpace(matches[1]), ToolInput: strings.TrimSpace(matches[2]), Log: output},
	}, nil, nil
}

func constructScratchPad(steps []schema.AgentStep) string {
	var scratchPad string
	if len(steps) > 0 {
		for _, step := range steps {
			scratchPad += step.Action.Log
			scratchPad += "\nObservation: " + step.Observation
		}
		scratchPad += "\n" + "Thought:"
	}

	return scratchPad
}
