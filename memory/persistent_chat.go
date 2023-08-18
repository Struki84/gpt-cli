package memory

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrDBConnection     = errors.New("can't connect to database")
	ErrDBMigration      = errors.New("can't migrate database")
	ErrMissingSessionID = errors.New("session id can not be empty")
)

type ChatHistory struct {
	ID           int       `gorm:"primary_key"`
	SessionID    string    `gorm:"type:varchar(256)"`
	BufferString string    `gorm:"type:text"`
	ChatHistory  *Messages `json:"chat_history" gorm:"type:jsonb;column:chat_history"`
}

type Messages []Message

type Message struct {
	Type    string `json:"type"`
	Content string `json:"text"`
}

type PersistentChatHistory struct {
	db        *gorm.DB
	records   *ChatHistory
	messages  []schema.ChatMessage
	sessionID string
}

var _ schema.ChatMessageHistory = &PersistentChatHistory{}

func NewPersistentChatHistory(dsn string) *PersistentChatHistory {

	history := &PersistentChatHistory{}

	gorm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	history.db = gorm

	err = history.db.AutoMigrate(ChatHistory{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return history
}

func (history *PersistentChatHistory) GetSessionID() string {
	return history.sessionID
}

func (history *PersistentChatHistory) SetSessionID(id string) {
	history.sessionID = id
}

func (history *PersistentChatHistory) Messages() ([]schema.ChatMessage, error) {
	if history.sessionID == "" {
		return []schema.ChatMessage{}, ErrMissingSessionID
	}

	err := history.db.Where(ChatHistory{SessionID: history.sessionID}).Find(&history.records).Error
	if err != nil {
		return nil, err
	}

	history.messages = []schema.ChatMessage{}

	if history.records.ChatHistory != nil {
		for i := range *history.records.ChatHistory {
			msg := (*history.records.ChatHistory)[i]

			if msg.Type == "human" {
				history.messages = append(history.messages, schema.HumanChatMessage{Content: msg.Content})
			}

			if msg.Type == "ai" {
				history.messages = append(history.messages, schema.AIChatMessage{Content: msg.Content})
			}
		}
	}

	return history.messages, nil
}

func (history *PersistentChatHistory) AddMessage(message schema.ChatMessage) error {
	if history.sessionID == "" {
		return ErrMissingSessionID
	}

	history.messages = append(history.messages, message)
	bufferString, err := schema.GetBufferString(history.messages, "Human", "AI")
	if err != nil {
		return err
	}

	history.records.SessionID = history.sessionID
	history.records.ChatHistory = history.loadNewMessages()
	history.records.BufferString = bufferString

	err = history.db.Save(&history.records).Error
	if err != nil {
		return err
	}

	return nil
}

func (history *PersistentChatHistory) AddAIMessage(message string) error {
	return history.AddMessage(schema.AIChatMessage{Content: message})
}

func (history *PersistentChatHistory) AddUserMessage(message string) error {
	return history.AddMessage(schema.HumanChatMessage{Content: message})
}

func (history *PersistentChatHistory) SetMessages(messages []schema.ChatMessage) error {
	if history.sessionID == "" {
		return ErrMissingSessionID
	}

	history.messages = messages
	bufferString, err := schema.GetBufferString(history.messages, "Human", "AI")
	if err != nil {
		return err
	}

	history.records.SessionID = history.sessionID
	history.records.ChatHistory = history.loadNewMessages()
	history.records.BufferString = bufferString

	return history.db.Save(&history.records).Error
}

func (history *PersistentChatHistory) Clear() error {
	history.messages = []schema.ChatMessage{}

	err := history.db.Where(ChatHistory{SessionID: history.sessionID}).Delete(&history.records).Error
	if err != nil {
		return err
	}

	return nil
}

func (history *PersistentChatHistory) loadNewMessages() *Messages {
	newMsgs := Messages{}
	for _, msg := range history.messages {
		newMsgs = append(newMsgs, Message{
			Type:    string(msg.GetType()),
			Content: msg.GetContent(),
		})
	}

	return &newMsgs
}

// Value implements the driver.Valuer interface, this method allows us to
// customize how we store the Message type in the database.
func (m Messages) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface, this method allows us to
// define how we convert the Message data from the database into our Message type.
func (m *Messages) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, m)
	}
	return errors.New("could not scan type into Message")
}
