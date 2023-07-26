package schema

import "github.com/tmc/langchaingo/schema"

type ChatMessageHistoryStore interface {
	
	// AddMessage method for storing a message in the DB store
	AddMessage(msg schema.ChatMessage) error

	// SetMessages method for replacing existing messages in the DB store
	SetMessages(msgs []schema.ChatMessage) error
	
	//GetMessages Convinience method for getting messages from db store
	GetMessages() ([]schema.ChatMessage, error)
	
	// ClearMessages method for clearing messages in the DB store
	ClearMessages() error
	
	// SetSessionId method for setting user session id
	SetSessionId(id string)

	// GetSessionId method for getting user session id
	GetSessionId() string
}