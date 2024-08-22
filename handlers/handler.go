package handlers

import (
	"fmt"
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
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	mes, err := models.MessageFromJSON(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling request body: %v", err),
			http.StatusBadRequest)
		return
	}
	err = models.ValidateMessage(mes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid message: %v", err), http.StatusBadRequest)
		return
	}
	err = h.st.Save(mes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Sorry! We failed to save the message: %v", err),
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	log.Println("Message saved", mes)
}

// Get возвращает статистику
func (h *Handler) Get(w http.ResponseWriter, _ *http.Request) {
	stats := h.st.GetStats()
	body, err := models.ToJSON(stats)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal stats: %v", err),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}
