//

package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/tmc/langchaingo/schema"
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
	history   *ChatHistory
	sessionID string
}

func NewPostgreAdapter() (*PostgreAdapter, error) {
	adapter := &PostgreAdapter{
		history: &ChatHistory{},
	}
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

func (adapter *PostgreAdapter) SetSessionId(id string) {
	adapter.sessionID = id
}

func (adapter *PostgreAdapter) GetSessionId() string {
	return adapter.sessionID
}

func (adapter *PostgreAdapter) SaveDBContext(id string, msgs []schema.ChatMessage, bufferString string) error {
	if adapter.sessionID == "" {
		return ErrMissingSessionID
	}

	newMsgs := Messages{}
	for _, msg := range msgs {
		newMsgs = append(newMsgs, Message{
			Type: string(msg.GetType()),
			Text: msg.GetText(),
		})
	}

	adapter.history.SessionID = adapter.sessionID
	adapter.history.ChatHistory = &newMsgs
	adapter.history.BufferString = bufferString

	err := adapter.gorm.Save(&adapter.history).Error
	if err != nil {
		return err
	}

	return nil
}

func (adapter *PostgreAdapter) LoadDBMemory(id string) ([]schema.ChatMessage, error) {
	// You implement your custom retrival logic here
	return nil, nil
}

func (adapter *PostgreAdapter) ClearDBContext(id string) error {
	// You implement your custom delete logic here
	return nil
}

// func main() {

// 	postgreAdapter, err := NewPostgreAdapter()
// 	persistentMemoryBuffer := memory.NewPersistentBuffer(postgreAdapter)
// 	persistentMemoryBuffer.DB.SetSessionId("USID-001")

// 	llm, err := openai.New()
// 	if err != nil {
// 		log.Print(err)
// 	}

// 	serpapi, err := serpapi.New()
// 	if err != nil {
// 		log.Print(err)
// 	}

// 	executor, err := agents.Initialize(
// 		llm,
// 		[]tools.Tool{serpapi},
// 		agents.ZeroShotReactDescription,
// 		agents.WithMemory(persistentMemoryBuffer),
// 		agents.WithMaxIterations(3),
// 	)

// 	if err != nil {
// 		log.Print(err)
// 	}

// 	input := "Who is the current CEO of Twitter?"
// 	answer, err := chains.Run(context.Background(), executor, input)
// 	if err != nil {
// 		log.Print(err)
// 		return
// 	}

// 	log.Print(answer)

// }

// An example of my custom DB adabter, a gorm based wrapper for a postgres database,
// implments the DBAdapter interface.
