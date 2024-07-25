package storage

import "message/models"

type Storage interface {
	Save(*models.Message) error
	GetStats() *models.Stats
}

func Save(st Storage, req *models.Message) error {
	return st.Save(req)
}

func GetStats(st Storage) *models.Stats {
	return st.GetStats()
}
