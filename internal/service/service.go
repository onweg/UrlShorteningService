package service

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/go-playground/validator/v10"
)

var data map[string]string = make(map[string]string)

const alphabetsKey = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const countLetterInId = 8

var validate *validator.Validate = validator.New()

type URLValidator struct {
	URL string `validate:"required,url"`
}

func ShorteningUrl(url string) (string, error) {
	if !isValidUrl(url) {
		return "", fmt.Errorf("%s invalid url", url)
	}

	randKey := getRandString()
	_, keyUsed := data[randKey]
	if keyUsed {
		return "", fmt.Errorf("%s key used", randKey)
	}

	data[randKey] = url
	log.Printf("Add new url: %s to key: %s\n", url, randKey)

	return randKey, nil
}

func GetOriginalUrl(shortUrl string) (string, error) {
	url := ""
	if !isValidId(shortUrl) {
		return "", fmt.Errorf("%s invalid id", shortUrl)
	}
	url, ok := data[shortUrl]
	if !ok {
		return "", fmt.Errorf("Not found url by id: %s", shortUrl)
	}
	return url, nil
}

func isValidUrl(url string) bool {
	val := URLValidator{URL: url}
	err := validate.Struct(val)
	return err == nil
}

func isValidId(id string) bool {
	for _, ch := range id {
		if !strings.Contains(alphabetsKey, string(ch)) {
			return false
		}
	}
	return true
}

func getRandString() string {
	res := make([]byte, countLetterInId)
	for i := 0; i < countLetterInId; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabetsKey))))
		res[i] = alphabetsKey[num.Int64()]
	}
	return string(res)
}
