package repository

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/go-playground/validator/v10"
)

const alphabetsKey = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const countLetterInId = 8

type URLValidator struct {
	URL string `validate:"required,url"`
}

type MemoryStore struct {
	data      map[string]string
	validator *validator.Validate
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:      make(map[string]string),
		validator: validator.New(),
	}
}

func (m MemoryStore) isValidUrl(url string) bool {
	val := URLValidator{URL: url}
	err := m.validator.Struct(val)
	return err == nil
}

func (m MemoryStore) isValidId(id string) bool {
	for _, ch := range id {
		if !strings.Contains(alphabetsKey, string(ch)) {
			return false
		}
	}
	return true
}

func (m *MemoryStore) Get(id string) (string, error) {
	url := ""
	if !m.isValidId(id) {
		return "", fmt.Errorf("%s invalid id", id)
	}
	url, ok := m.data[id]
	if !ok {
		return "", fmt.Errorf("Not found url by id: %s", id)
	}
	return url, nil
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

func getRandString() string {
	res := make([]byte, countLetterInId)
	for i := 0; i < countLetterInId; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabetsKey))))
		res[i] = alphabetsKey[num.Int64()]
	}
	return string(res)
}
