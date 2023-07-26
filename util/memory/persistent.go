package memory

import (
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

type DBAdapter interface {
	LoadDBMemory(id string) ([]schema.ChatMessage, error)
	SaveDBContext(id string, msgs []schema.ChatMessage, bufferString string) error
	ClearDBContext(id string) error
	SetSessionId(id string)
	GetSessionId() string
}

type PersistentBuffer struct {
	ChatHistory *memory.ChatMessageHistory
	DB          DBAdapter

	ReturnMessages bool
	InputKey       string
	OutputKey      string
	HumanPrefix    string
	AIPrefix       string
	MemoryKey      string
}

var _ schema.Memory = PersistentBuffer{}

func NewPersistentBuffer(dbAdapter DBAdapter) *PersistentBuffer {
	buffer := PersistentBuffer{
		ChatHistory: memory.NewChatMessageHistory(),
		DB:          dbAdapter,

		ReturnMessages: false,
		InputKey:       "",
		OutputKey:      "",
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		MemoryKey:      "history",
	}

	return &buffer
}

func (buffer PersistentBuffer) MemoryVariables() []string {
	return []string{buffer.MemoryKey}
}

func (buffer PersistentBuffer) LoadMemoryVariables(inputs map[string]any) (map[string]any, error) {
	sessionID := buffer.DB.GetSessionId()
	msgs, err := buffer.DB.LoadDBMemory(sessionID)
	if err != nil {
		return nil, err
	}

	buffer.ChatHistory = memory.NewChatMessageHistory(
		memory.WithPreviousMessages(msgs),
	)

	msgs, err = buffer.ChatHistory.Messages()

	if buffer.ReturnMessages {
		return map[string]any{
			buffer.MemoryKey: msgs,
		}, nil
	}

	bufferString, err := schema.GetBufferString(msgs, buffer.HumanPrefix, buffer.AIPrefix)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		buffer.MemoryKey: bufferString,
	}, nil
}

func (buffer PersistentBuffer) SaveContext(inputs map[string]any, outputs map[string]any) error {
	sessionID := buffer.DB.GetSessionId()
	userInputValue, err := getInputValue(inputs, buffer.InputKey)
	if err != nil {
		return err
	}

	buffer.ChatHistory.AddUserMessage(userInputValue)

	aiOutPutValue, err := getInputValue(outputs, buffer.OutputKey)
	if err != nil {
		return err
	}

	buffer.ChatHistory.AddAIMessage(aiOutPutValue)

	msgs, err := buffer.ChatHistory.Messages()
	if err != nil {
		return err
	}

	bufferString, err := schema.GetBufferString(msgs, buffer.HumanPrefix, buffer.AIPrefix)
	if err != nil {
		return err
	}

	err = buffer.DB.SaveDBContext(sessionID, msgs, bufferString)
	if err != nil {
		return err
	}

	return nil
}

func (buffer PersistentBuffer) GetMemoryKey() string {
	return buffer.MemoryKey
}

func (buffer PersistentBuffer) Clear() error {
	sessionID := buffer.DB.GetSessionId()

	err := buffer.DB.ClearDBContext(sessionID)
	if err != nil {
		return err
	}

	buffer.ChatHistory.Clear()

	return nil
}
