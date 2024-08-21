package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

// ToJSON парсит в JSON
func ToJSON(v any) ([]byte, error) {
	return json.Marshal(v)
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

// ValidateMessage проверяет структуру на валидность
// и убирает символы пробелов и табуляции в начале и конце
// в полях Author и Text
func ValidateMessage(mes *Message) error {
	if mes == nil {
		return fmt.Errorf("author or text is empty")
	}
	mes.Author = strings.Trim(mes.Author, " \t\r")
	mes.Text = strings.Trim(mes.Text, " \t\r")
	if mes.Author == "" || mes.Text == "" {
		return fmt.Errorf("author or text is empty/invalid")
	}
	return nil
}

// String возвращает строковое представление сообщения
func (m *Message) String() string {
	return m.Author + ": " + m.Text
}
