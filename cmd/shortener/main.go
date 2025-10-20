package main

import (
	"net/http"

	"github.com/onweg/UrlShorteningService/internal/handler"
	"github.com/onweg/UrlShorteningService/internal/service"
	"github.com/onweg/UrlShorteningService/internal/model"
)

func main() {
	storage := model.NewMemoryStorage()
    urlService := service.NewURLService(storage)
    h := handler.NewHandler(urlService)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HandleRequest())
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
