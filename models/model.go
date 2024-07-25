package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Message struct {
	ID      int
	Author  string `json:"author"`
	Text    string `json:"text"`
	Created time.Time
	Sent    bool
}

type Stats struct {
	Counter int
}

func MessageFromJson(body []byte) (*Message, error) {
	var mes Message

	err := json.Unmarshal(body, &mes)
	//TODO add validation function
	if err != nil {
		fmt.Println("Error unmarshalling body:", err)
		return nil, err
	}
	return &mes, nil
}

func StatsToJSON(stats *Stats) ([]byte, error) {
	body, err := json.Marshal(stats)
	if err != nil {
		log.Fatal("Marshalling body error:", err)
		return nil, err
	}
	return body, nil
}

func (m *Message) String() string {
	return m.Author + ": " + m.Text
}
