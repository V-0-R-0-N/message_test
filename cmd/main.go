package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"message/handlers"
	storage "message/storage/simple"
	"net/http"
)

var addr string = ":8080"

type Message struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

func main() {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	st := storage.NewSimpleStorage()
	h := handlers.NewHandler(st)
	r.Post("/send", h.Save)
	r.Get("/stats", h.Get)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
