package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/onweg/UrlShorteningService/internal/service"
)

func HandleRequest() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			handleGet(res, req)
		case http.MethodPost:
			handlePost(res, req)
		default:
			res.Write([]byte("the service supports only GET and POST requests"))
		}
	}
}

func handlePost(res http.ResponseWriter, req *http.Request) {
	resBody, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		// не смог прочитать лучше, вернуть 500 ошибку
		res.WriteHeader(http.StatusBadRequest)
		log.Printf("couldn't read the request: %s\n", err.Error())
		return
	}

	shortUrl, err := service.ShorteningUrl(string(resBody))
	if err != nil {
		// не смог сделать ключ, вернуть 500 ошибку
		res.WriteHeader(http.StatusBadRequest)
		log.Printf("error: %s\n", err.Error())
		return
	}

	outResult := "http://" + req.Host + "/" + shortUrl
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(outResult))
}

func handleGet(res http.ResponseWriter, req *http.Request) {
	keyFindUrl := req.RequestURI[1:]
	if keyFindUrl == "" {
		res.Write([]byte("Добро пожаловать!"))
		return
	}

	url, err := service.GetOriginalUrl(keyFindUrl)
	log.Printf("Get url: %s to key: %s\n", url, keyFindUrl)
	if err != nil {
		// некорректный id не нашли по id нужный url, вернуть 400
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Redirect to url: %s\n", url)
	res.Header().Set("Location", url)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
