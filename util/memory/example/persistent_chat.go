package main

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gpt/util/memory"
	"gpt/util/schema"
	"log"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	mem "github.com/tmc/langchaingo/memory"
	sch "github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/serpapi"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrDBConnection = errors.New("can't connect to database")
var ErrDBMigration = errors.New("can't migrate database")
var ErrMissingSessionID = errors.New("session id can not be empty")

type ChatHistory struct {
	ID           int       `gorm:"primary_key"`
	SessionID    string    `gorm:"type:varchar(256)"`
	BufferString string    `gorm:"type:text"`
	ChatHistory  *Messages `json:"chat_history" gorm:"type:jsonb;column:chat_history"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Messages []Message

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

type PostgreAdapter struct {
	gorm      *gorm.DB
	sessionID string
	history   *ChatHistory
}

var _ schema.ChatMessageHistoryStore = &PostgreAdapter{}

func NewPostgreAdapter() (*PostgreAdapter, error) {
	adapter := &PostgreAdapter{}
	
	dsn := ""
	
	gorm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, ErrDBConnection
	}

	adapter.gorm = gorm

	err = adapter.gorm.AutoMigrate(ChatHistory{})
	if err != nil {
		return nil, ErrDBMigration
	}

	return adapter, nil
}

func (adapter *PostgreAdapter) SetSessionID(id string) {
	adapter.sessionID = id
}

func (adapter *PostgreAdapter) GetSessionID() string {
	return adapter.sessionID
}

func (adapter *PostgreAdapter) AddMessage(msg sch.ChatMessage) error {
	if adapter.sessionID == "" {
		return ErrMissingSessionID
	}

	msgs, err := adapter.GetMessages()
	if err != nil {
		return err
	}

	msgs = append(msgs, msg)
	err = adapter.SetMessages(msgs)
	if err != nil {
		return err
	}

	return nil
}

func (adapter *PostgreAdapter) SetMessages(msgs []sch.ChatMessage) error {
	if adapter.sessionID == "" {
		return ErrMissingSessionID
	}

	newMsgs := Messages{}
	for _, msg := range msgs {
		newMsgs = append(newMsgs, Message{
			Type: string(msg.GetType()),
			Text: msg.GetContent(),
		})
	}

	adapter.history.SessionID = adapter.sessionID
	adapter.history.ChatHistory = &newMsgs
	
	err := adapter.gorm.Save(&adapter.history).Error
	if err != nil {
		return err
	}

	return nil
}

func (adapter *PostgreAdapter) GetMessages() ([]sch.ChatMessage, error) {
	if adapter.sessionID == "" {
		return nil, ErrMissingSessionID
	}

	err := adapter.gorm.Where(ChatHistory{SessionID: adapter.sessionID}).Find(&adapter.history).Error
	if err != nil {
		return nil, err
	}

	msgs := []sch.ChatMessage{}
	if adapter.history.ChatHistory != nil {
		for i := range *adapter.history.ChatHistory {
			msg := (*adapter.history.ChatHistory)[i]

			if msg.Type == "human" {
				msgs = append(msgs, sch.HumanChatMessage{Content: msg.Text})
			}

			if msg.Type == "ai" {
				msgs = append(msgs, sch.AIChatMessage{Content: msg.Text})
			}
		}
	}

	return msgs, nil
}

func (adapter *PostgreAdapter) ClearMessages() error {
	err := adapter.gorm.Where(ChatHistory{SessionID: adapter.sessionID}).Delete(&adapter.history).Error
	if err != nil {
		return err
	}
	return nil
}

func main() {

	postgreAdapter, err := NewPostgreAdapter()
	if err != nil {
		log.Print(err)
	}

	chatHistory := memory.NewPersistentChatMessageHistory(memory.WithDBStore(postgreAdapter))
	memoryBuffer := mem.NewConversationBuffer(mem.WithChatHistory(chatHistory))

	llm, err := openai.New()
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	serpapi, err := serpapi.New()
	if err != nil {
		log.Print(err)
	}

	iterations := 3

	executor, err := agents.Initialize(
		llm,
		[]tools.Tool{serpapi},
		agents.ZeroShotReactDescription,
		agents.WithMemory(memoryBuffer),
		agents.WithMaxIterations(iterations),
	)
	if err != nil {
		log.Print(err)
	}

	input := "Who is the current CEO of Twitter?"
	answer, err := chains.Run(context.Background(), executor, input)
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(answer)
}





