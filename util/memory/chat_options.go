package memory

import (
	"gpt/util/schema"

	sch "github.com/tmc/langchaingo/schema"
)

// ChatMessageHistoryOption is a function for creating new chat message history
// with other then the default values.
type PersistentChatMessageHistoryOption func(m *PersistentChatMessageHistory)

// WithPreviousMessages is an option for NewChatMessageHistory for adding
// previous messages to the history.
func WithPreviousMessages(previousMessages []sch.ChatMessage) PersistentChatMessageHistoryOption {
	return func(m *PersistentChatMessageHistory) {
		m.messages = append(m.messages, previousMessages...)
	}
}

func WithDBStore(dbstore schema.ChatMessageHistoryStore) PersistentChatMessageHistoryOption {
	return func(m *PersistentChatMessageHistory) {
		m.store = dbstore
	}
}

func applyChatOptions(options ...PersistentChatMessageHistoryOption) *PersistentChatMessageHistory {
	h := &PersistentChatMessageHistory{
		messages: make([]sch.ChatMessage, 0),
	}

	for _, option := range options {
		option(h)
	}

	return h
}