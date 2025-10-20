package model

import "fmt"

type Storage interface {
	Save(key, url string) error
	Get(key string) (string, error)
}

type MemoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (s *MemoryStorage) Save(key, url string) error {
	if _, ok := s.data[key]; ok {
		return fmt.Errorf("%s key used", key)
	}
	s.data[key] = url
	return nil
}

func (s *MemoryStorage) Get(key string) (string, error) {
	url, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("Not found url by id: %s", url)
	}
	return url, nil
}
