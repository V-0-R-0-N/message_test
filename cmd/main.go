package main

import (
	"log"
	"message/internal/router"
	storage "message/storage/postgres"
	"net/http"
)

var addr string = ":8080"

func main() {

	//st := storage.NewSimple()
	db := storage.NewDB()
	defer db.Pool.Close()
	r := router.InitNew(db)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
