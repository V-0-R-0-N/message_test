package main

import (
	"log"
	"message/internal/router"
	storage "message/storage/simple"
	"net/http"
)

var addr string = ":8080"

func main() {

	st := storage.NewSimple()

	r := router.InitNew(st)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
