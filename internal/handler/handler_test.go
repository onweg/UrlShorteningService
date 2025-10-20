package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/onweg/UrlShorteningService/internal/handler"
	"github.com/onweg/UrlShorteningService/internal/model"
	"github.com/onweg/UrlShorteningService/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLHandler_handleGet(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    string
	}
	tests := []struct {
		name    string // description of this test case
		request string
		storage map[string]string
		want    want
	}{
		{
			name:    "success get request 1",
			request: "/u1",
			storage: map[string]string{
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code:        307,
				contentType: "text/plain",
				response:    ``,
			},
		},
		{
			name:    "success get request 2",
			request: "/u2",
			storage: map[string]string{
				"u2": "https://www.google.com/",
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code:        307,
				contentType: "text/plain",
				response:    ``,
			},
		},
		{
			name:    "faile get request 1",
			request: "/invalidID",
			storage: map[string]string{
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
				response:    ``,
			},
		},
		{
			name:    "faile get request 2",
			request: "/u1",
			storage: map[string]string{
				// empty
			},
			want: want{
				code:        400,
				contentType: "text/plain",
				response:    ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &model.MemoryStorage{
				Data: tt.storage,
			}
			service := service.NewURLService(storage)
			handler := handler.NewHandler(service)

			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.HandleRequest())
			h(w, r)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}

func TestURLHandler_handlePost(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    string
	}
	tests := []struct {
		name    string // description of this test case
		request string
		body    string
		storage map[string]string
		want    want
	}{
		{
			name:    "success post request 1",
			request: "/",
			body:    `https://practicum.yandex.ru/`,
			storage: map[string]string{
				// empty
			},
			want: want{
				code:        201,
				contentType: "text/plain",
				response:    `^http://localhost:8080/[A-Za-z0-9]{8}$`,
			},
		},
		{
			name:    "success post request 2",
			request: "/",
			body:    `https://www.google.com/`,
			storage: map[string]string{
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code:        201,
				contentType: "text/plain",
				response:    `^http://localhost:8080/[A-Za-z0-9]{8}$`,
			},
		},
		{
			name:    "success post request 3",
			request: "/",
			body:    `https://www.google.com/`,
			storage: map[string]string{
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code:        201,
				contentType: "text/plain",
				response:    `^http://localhost:8080/[A-Za-z0-9]{8}$`,
			},
		},
		{
			name:    "faile post request 1",
			request: "/",
			body:    `https/www.google.com/`,
			storage: map[string]string{
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
				response:    ``,
			},
		},
		{
			name:    "faile post request 2",
			request: "/",
			body:    ``,
			storage: map[string]string{
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
				response:    ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &model.MemoryStorage{
				Data: tt.storage,
			}
			service := service.NewURLService(storage)
			handler := handler.NewHandler(service)

			r := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.HandleRequest())
			r.Host = "localhost:8080"
			h(w, r)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			if result.StatusCode == http.StatusCreated {
				resBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				assert.Regexp(t, tt.want.response, string(resBody))
			}
		})
	}
}
