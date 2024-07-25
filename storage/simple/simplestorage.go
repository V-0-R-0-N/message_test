package storage

import (
	"message/models"
	"strconv"
	"strings"
	"time"
)

type Simple struct {
	Counter int
	Data    []*data
}

type data struct {
	ID      int
	Author  string
	Text    string
	Created time.Time
}

func (s *Simple) Save(req *models.Message) error {
	s.Data = append(s.Data, &data{
		ID:      s.Counter,
		Author:  req.Author,
		Text:    req.Text,
		Created: time.Now(),
	})
	s.Counter++
	return nil
}

func (s *Simple) GetStats() *models.Stats {
	return &models.Stats{
		Counter: s.Counter,
	}
}

func NewSimpleStorage() *Simple {
	return &Simple{}
}
func (s *Simple) String() string {
	var str strings.Builder

	for _, v := range s.Data {
		str.WriteString("ID: " + strconv.Itoa(v.ID))
		str.WriteString(" ")
		str.WriteString("Author: " + v.Author)
		str.WriteString(" ")
		str.WriteString("Text: " + v.Text)
		str.WriteString(" ")
		str.WriteString("Created: " + v.Created.String())
		str.WriteString("\n")
	}
	return str.String()
}
