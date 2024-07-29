package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"message/handlers"
	"message/storage"
)

// InitNew инициализирует новый роутер
func InitNew(st storage.Storage) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handlers.NewHandler(st)
	r.Post("/send", h.Save)
	r.Get("/stats", h.Get)

	return r
}
