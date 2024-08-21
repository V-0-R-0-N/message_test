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
