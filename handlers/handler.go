package handlers

import (
	"io"
	"log"
	"message/models"
	"message/storage"
	"net/http"
)

// Handler обработчик запросов с хранилищем
type Handler struct {
	st storage.Storage
}

// NewHandler создает новый обработчик
func NewHandler(st storage.Storage) *Handler {
	return &Handler{
		st: st,
	}
}

// Save сохраняет сообщение
func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}
	mes, err := models.MessageFromJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error or invalid request\n"))
		return
	}
	err = models.ValidateMessage(mes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		log.Println(err)
		return
	}
	storage.Save(h.st, mes)
	w.WriteHeader(http.StatusCreated)

	log.Println("Message saved", mes)
}

// Get возвращает статистику
func (h *Handler) Get(w http.ResponseWriter, _ *http.Request) {
	stats := storage.GetStats(h.st)
	body, err := models.StatsToJSON(stats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}
