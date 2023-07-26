package memory

import (
	"gpt/util/schema"

	sch "github.com/tmc/langchaingo/schema"
)

type PersistentChatMessageHistory struct {
	messages []sch.ChatMessage
	store 	 schema.ChatMessageHistoryStore
}

var _ sch.ChatMessageHistory = &PersistentChatMessageHistory{}

func NewPersistentChatMessageHistory(options ...PersistentChatMessageHistoryOption,) *PersistentChatMessageHistory {
	return applyChatOptions(options...)
}

func (h *PersistentChatMessageHistory)Messages() ([]sch.ChatMessage, error) {
	msgs, err := h.store.GetMessages()
	if err != nil {
		return nil, err
	}
	h.messages = msgs
	return h.messages, nil
}

func (h *PersistentChatMessageHistory)AddAIMessage(text string) error {
	msg := sch.AIChatMessage{Content: text}
	
	err := h.store.AddMessage(msg)
	if err != nil {
		return err
	}

	h.messages = append(h.messages, msg)
	return nil
}

func (h *PersistentChatMessageHistory)AddUserMessage(text string) error {
	msg := sch.HumanChatMessage{Content: text}
	
	err := h.store.AddMessage(msg)
	if err != nil {
		return err
	}

	h.messages = append(h.messages, msg)
	return nil
}

func (h *PersistentChatMessageHistory)AddMessage(message sch.ChatMessage) error {
	err := h.store.AddMessage(message)
	
	if err != nil {
		return err
	}
	
	h.messages = append(h.messages, message)
	return nil
}

func (h *PersistentChatMessageHistory) SetMessages(msgs []sch.ChatMessage) error {
	err := h.store.SetMessages(msgs)
	if err != nil {
		return err
	}
	h.messages = msgs
	return nil
	
}

func (h *PersistentChatMessageHistory) Clear() error {
	err := h.store.ClearMessages()
	if err != nil {
		return err
	}
	
	h.messages = make([]sch.ChatMessage, 0)
	return nil
}

