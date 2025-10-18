package main

import (
	"io"
	"net/http"
	"fmt"
	"crypto/rand"
	"math/big"
	"strings"
	"log"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type URLValidator struct {
	URL string `validate:"required,url"`
}

const alphabetsKey = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const countLetterInId = 8

type URLStore interface {
	Save(url string)(string, error)
	Get(id string)(string, error)
}

type MemoryStore struct {
	data map[string]string
}

func initValidate() {
    validate = validator.New()
}

func (m *MemoryStore) Save(url string) (string, error) {
	if !m.isValidUrl(url) {
		return "", fmt.Errorf("%s invalid url", url)
	}

	randKey := getRandString()
	_, keyUsed := m.data[randKey]
	if keyUsed {
		return "", fmt.Errorf("%s key used", randKey)
	}

	m.data[randKey] = url
	log.Printf("Add new url: %s to key: %s\n", url, randKey)

	return randKey, nil
}

func (m MemoryStore) isValidUrl(url string) bool {
	val := URLValidator {URL: url}
	err := validate.Struct(val)
	if err != nil {
		return false
	}
	return true
}

func (m *MemoryStore) Get(id string) (string, error) {
	url := ""
	if !m.isValidId(id) {
		return "", fmt.Errorf("%s invalid id", id)
	}

	url = m.data[id]
	return url, nil
}

func (m MemoryStore) isValidId(id string) bool {
	for _, ch := range id {
		if strings.Contains(alphabetsKey, string(ch)) == false {
			return false
		}
	} 
	return true
}

type Server struct {
	store URLStore
}

func (s Server) handlePost(res http.ResponseWriter, req *http.Request){
	resBody, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		// не смог прочитать лучше, вернуть 500 ошибку
		res.WriteHeader(http.StatusBadRequest)
		log.Printf("couldn't read the request: %s\n", err.Error())
		return
	}

	newId, err := s.store.Save(string(resBody))
	if err != nil {
		// не смог сделать ключ, вернуть 500 ошибку
		res.WriteHeader(http.StatusBadRequest)
		log.Printf("error: %s\n", err.Error())
		return
	}

	outResult := "http://" + req.Host + "/" + newId
	res.Header().Set("Content-Type",  "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(outResult))
}

func (s Server) handleGet(res http.ResponseWriter, req *http.Request){
	keyFindUrl := req.RequestURI[1:]
	if keyFindUrl == "" {
		res.Write([]byte("Добро пожаловать!"))
		return
	}

	url, err := s.store.Get(keyFindUrl)
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

func getRandString() string {
	res := make([]byte, countLetterInId)
	for i := 0; i < countLetterInId; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabetsKey)))) 
		res[i] = alphabetsKey[num.Int64()]
	}
	return string(res)
}

func handleRequest(s Server) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
			case http.MethodGet:
				s.handleGet(res, req)
			case http.MethodPost:
				s.handlePost(res, req)
			default:
				res.Write([]byte("the service supports only GET and POST requests"))
		}
	}
}

func main() {
	initValidate()
	store := &MemoryStore{data: make(map[string]string)}
	server := Server{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRequest(server))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
