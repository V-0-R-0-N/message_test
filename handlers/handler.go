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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}
	err = models.ValidateMessage(mes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		log.Println(err)
		return
	}
	err = h.st.Save(mes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry! We failed to save the message\n"))
		log.Println(err)
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}
