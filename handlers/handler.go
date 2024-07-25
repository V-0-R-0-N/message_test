package handlers

import (
	"fmt"
	"io"
	"log"
	"message/models"
	"message/storage"
	"net/http"
)

type Handler struct {
	st storage.Storage
}

func NewHandler(st storage.Storage) *Handler {
	return &Handler{
		st: st,
	}
}

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}
	mes, err := models.MessageFromJson(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}
	storage.Save(h.st, mes)
	w.WriteHeader(http.StatusCreated)

	fmt.Println(mes)
	fmt.Println(h.st)
}

func (h *Handler) Get(w http.ResponseWriter, _ *http.Request) {
	stats := storage.GetStats(h.st)
	body, err := models.StatsToJSON(stats)
	if err != nil {
		log.Fatal("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}
