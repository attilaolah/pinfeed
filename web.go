package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", pinFeed)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
