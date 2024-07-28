package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Message структура сообщения
type Message struct {
	ID      int
	Author  string `json:"author"`
	Text    string `json:"text"`
	Created time.Time
	Sent    bool
}

// Stats структура статистики
type Stats struct {
	Total int
	Sent  int
}

// MessageFromJSON парсит JSON в структуру Message
func MessageFromJSON(body []byte) (*Message, error) {
	var mes Message

	err := json.Unmarshal(body, &mes)
	if err != nil {
		fmt.Println("Error unmarshalling body:", err)
		return nil, err
	}
	return &mes, nil
}

// StatsToJSON парсит структуру Stats в JSON
func StatsToJSON(stats *Stats) ([]byte, error) {
	body, err := json.Marshal(stats)
	if err != nil {
		log.Fatal("Marshalling body error:", err)
		return nil, err
	}
	return body, nil
}

// MessageToJSON парсит структуру Message в JSON
func MessageToJSON(mes *Message) ([]byte, error) {
	body, err := json.Marshal(mes)
	if err != nil {
		log.Fatal("Marshalling body error:", err)
		return nil, err
	}
	return body, nil
}

// MessageToJSONForKafka парсит структуру Message в JSON для Kafka
func MessageToJSONForKafka(mes *Message) ([]byte, error) {
	body, err := json.Marshal(struct {
		ID      int       `json:"id"`
		Author  string    `json:"author"`
		Text    string    `json:"text"`
		Created time.Time `json:"created"`
	}{
		ID:      mes.ID,
		Author:  mes.Author,
		Text:    mes.Text,
		Created: mes.Created,
	})
	if err != nil {
		log.Fatal("Marshalling body error:", err)
		return nil, err
	}
	return body, nil
}

// ValidateMessage проверяет сообщение на валидность
func ValidateMessage(mes *Message) error {
	if mes.Author == "" || mes.Text == "" {
		return fmt.Errorf("author or text is empty")
	}
	return nil
}

// String возвращает строковое представление сообщения
func (m *Message) String() string {
	return m.Author + ": " + m.Text
}
