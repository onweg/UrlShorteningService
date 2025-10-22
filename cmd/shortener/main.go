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
	r := handler.NewRouter(h)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
