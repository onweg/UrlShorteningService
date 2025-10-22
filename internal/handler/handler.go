package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/onweg/UrlShorteningService/internal/service"
)

type Handler struct {
	urlService *service.URLService
}

func NewRouter(h *Handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/{shortUrl}", h.getOriginalHandle)
	r.Post("/", h.handlePost)
	return r
}

func NewHandler(urlService *service.URLService) *Handler {
	return &Handler{urlService: urlService}
}

func (h *Handler) handlePost(res http.ResponseWriter, req *http.Request) {
	resBody, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		// не смог прочитать лучше, вернуть 500 ошибку
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusBadRequest)
		log.Printf("couldn't read the request: %s\n", err.Error())
		return
	}

	shortUrl, err := h.urlService.ShorteningUrl(string(resBody))
	if err != nil {
		// не смог сделать ключ, вернуть 500 ошибку
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusBadRequest)
		log.Printf("error: %s\n", err.Error())
		return
	}

	outResult := "http://" + req.Host + "/" + shortUrl
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(outResult))
}

func (h *Handler) getOriginalHandle(res http.ResponseWriter, req *http.Request) {
	shortUrl := chi.URLParam(req, "shortUrl")
	if shortUrl == "" {
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte("Welcom!"))
		return
	}

	url, err := h.urlService.GetOriginalUrl(shortUrl)
	log.Printf("Get url: %s to key: %s\n", url, shortUrl)
	if err != nil {
		// некорректный id не нашли по id нужный url, вернуть 400
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Redirect to url: %s\n", url)
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", url)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
