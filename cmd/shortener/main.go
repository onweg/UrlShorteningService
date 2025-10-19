package main

import (
	"net/http"

	"github.com/onweg/UrlShorteningService/internal/handler"
	"github.com/onweg/UrlShorteningService/internal/repository"
	"github.com/onweg/UrlShorteningService/internal/service"
)

type URLStore interface {
	Save(url string) (string, error)
	Get(id string) (string, error)
}

type Server struct {
	store URLStore
}

func main() {
	store := repository.NewMemoryStore()
	svc := service.NewUrlService(store)
	h := handler.NewURLHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HandleRequest())
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
