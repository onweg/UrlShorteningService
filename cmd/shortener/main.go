package main

import (
	"net/http"

	"github.com/onweg/UrlShorteningService/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HandleRequest())
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
