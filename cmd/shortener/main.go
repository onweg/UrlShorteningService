package main

import (
	"io"
	"net/http"
	"fmt"
	"math/rand"
	"time"

)

var URL_API string = "http://localhost:8080"

func getRandString() string {
	res := ""
	for i := 0; i < 8; i++ {
		lowLet := rand.Int() % 2
		if lowLet == 1 {
			res += string(byte(rand.Int() % 26 + 97))
		} else {
			res += string(byte(rand.Int() % 26 + 65))
		}
	}
	return res
}

func proccesingServices(m map [string]string) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
			case http.MethodGet:
				keyFindUrl := req.RequestURI[1:]
				if keyFindUrl == "" {
					res.Write([]byte("Добро пожаловать!"))
					break
				}
				url := m[keyFindUrl]
				fmt.Printf("Get url: %s to key: %s\n", url, keyFindUrl)
				if url == "" {
					res.WriteHeader(http.StatusBadRequest)
					break
				}
				fmt.Printf("Redirect to url: %s\n", url)
				res.Header().Set("Location", url)
				http.Redirect(res, req, url, http.StatusTemporaryRedirect)
			case http.MethodPost:
				resBody, err := io.ReadAll(req.Body)
				if err != nil {
					panic(err)
				}
				randKey := getRandString()
				m[randKey] = string(resBody)
				fmt.Printf("Add new url: %s to key: %s\n", string(resBody), randKey)
				req.URL.Scheme = "http"				
				outResult := req.URL.Scheme + "://" + req.Host + "/" + randKey

				res.Header().Set("Content-Type",  "text/plain")
				res.Header().Set("Content-Length", fmt.Sprintf("%d", len(outResult)))
				res.WriteHeader(http.StatusCreated)
				
				res.Write([]byte(outResult))
			default:
				res.Write([]byte("the service supports only GET and POST requests"))
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	shortToUrl := make(map [string]string)
	mux := http.NewServeMux()
	mux.HandleFunc("/", proccesingServices(shortToUrl))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
