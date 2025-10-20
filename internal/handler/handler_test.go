package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/onweg/UrlShorteningService/internal/service"
)

func TestURLHandler_handleGet(t *testing.T) {
	type want struct {
		code string 
		contentType string 
		response string
	}
	tests := []struct {
		name string // description of this test case
		request string
		store map[string]string
		want want
	}{
		{
			name: "success get request",
			request: "/u1",
			store: map[string]string{
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code: "307",
				contentType: "text/plain",
				response: ``,
			},
		},
		{
			name: "success get request",
			request: "/u2",
			store: map[string]string{
				"u2": "https://www.google.com/",
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code: "307",
				contentType: "text/plain",
				response: ``,
			},
		},
		{
			name: "success get request",
			request: "/invalidID",
			store: map[string]string{
				"u1": "https://practicum.yandex.ru/",
			},
			want: want{
				code: "400",
				contentType: "text/plain",
				response: ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Hand)
		})
	}
}

// func TestURLHandler_HandleRequest(t *testing.T) {
// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for receiver constructor.
// 		svc  *service.URLService
// 		want func(res http.ResponseWriter, req *http.Request)
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := NewURLHandler(tt.svc)
// 			got := h.HandleRequest()
// 			// TODO: update the condition below to compare got with tt.want.
// 			if true {
// 				t.Errorf("HandleRequest() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
