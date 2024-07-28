package main

import (
	"context"
	"log"
	"message/internal/kafka"
	"message/internal/router"
	"message/internal/worker"
	storage "message/storage/postgres"
	"net/http"
)

var addr string = ":8080"

func main() {

	//st := storage.NewSimple()
	db := storage.NewDB()
	defer db.Pool.Close()
	ctx := context.Background()
	producer := kafka.InitProducer()
	defer (*producer).Close()
	go worker.New(ctx, db, producer)
	r := router.InitNew(db)
	log.Println("Server is running on", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
