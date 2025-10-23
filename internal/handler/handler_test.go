package handler_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/onweg/UrlShorteningService/internal/handler"
	"github.com/onweg/UrlShorteningService/internal/model"
	"github.com/onweg/UrlShorteningService/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func testRequest(t *testing.T, ts *httptest.Server, method, path, body string) (*http.Response, string) {

	var req *http.Request
	var err error

	if method == http.MethodGet {
		req, err = http.NewRequest(method, ts.URL+path, nil)
	} else if method == http.MethodPost {
		req, err = http.NewRequest(method, ts.URL+path, bytes.NewBufferString(body))
	}
	require.NoError(t, err)

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBosy, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBosy)
}

func TestRouterGetMethod(t *testing.T) {
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
	for _, v := range tests {
		storage := &model.MemoryStorage{Data: v.storage}
		urlService := service.NewURLService(storage)
		h := handler.NewHandler(urlService)
		r := handler.NewRouter(h)
		ts := httptest.NewServer(r)
		defer ts.Close()

		resp, get := testRequest(t, ts, http.MethodGet, v.request, "")
		assert.Equal(t, resp.StatusCode, v.want.code)
		assert.Regexp(t, v.want.response, string(get))
	}
}

func TestRouterPostMethod(t *testing.T) {
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
				response:    `^http://(localhost|127\.0\.0\.1):\d+/[A-Za-z0-9]{8}$`,
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
				response:    `^http://(localhost|127\.0\.0\.1):\d+/[A-Za-z0-9]{8}$`,
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
				response:    `^http://(localhost|127\.0\.0\.1):\d+/[A-Za-z0-9]{8}$`,
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
	for _, v := range tests {
		storage := &model.MemoryStorage{Data: v.storage}
		urlService := service.NewURLService(storage)
		h := handler.NewHandler(urlService)
		r := handler.NewRouter(h)
		ts := httptest.NewServer(r)
		defer ts.Close()

		resp, get := testRequest(t, ts, http.MethodPost, v.request, v.body)
		assert.Equal(t, resp.StatusCode, v.want.code)
		assert.Regexp(t, v.want.response, string(get))
	}
}
