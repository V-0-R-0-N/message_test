package storage

import (
	"context"
	"message/models"
)

// Storage интерфейс для хранилища
type Storage interface {
	Save(*models.Message) error
	GetStats() *models.Stats
	NeedSent() []*models.Message
	ChangeStatusSent(context.Context, int) error
}

// Save сохраняет сообщение
func Save(st Storage, req *models.Message) error {
	return st.Save(req)
}

// GetStats возвращает статистику принятых и отправленных сообщений
func GetStats(st Storage) *models.Stats {
	return st.GetStats()
}

// NeedSent возвращает сообщения, которые нужно отправить
func NeedSent(st Storage) []*models.Message {
	return st.NeedSent()
}

// ChangeStatusSent изменяет статус отправки сообщения
func ChangeStatusSent(ctx context.Context, st Storage, id int) error {
	return st.ChangeStatusSent(ctx, id)
}
