package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
)

var addr string = ":8080"

type Message struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

func simpleGet(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)

	fmt.Println("simple request")
}

func simplePost(w http.ResponseWriter, r *http.Request) {

	mes := Message{}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}
	err = json.Unmarshal(body, &mes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{"saved_status" : "ok"}`))
	w.WriteHeader(http.StatusOK)
}

func main() {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/send", simplePost)
	r.Get("/stats", simpleGet)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
