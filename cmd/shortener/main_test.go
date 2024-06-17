package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/config"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	links := storage.NewStorage()
	config := config.NewConfig()
	h := &handlers.Handler{
		Storage: links,
		Config:  config,
	}

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid URL",
			body:       "http://www.ya.ru/",
			wantStatus: http.StatusCreated,
			wantBody:   "http://localhost:8080/6f782b56",
		},
		{
			name:       "Invalid URL",
			body:       "",
			wantStatus: http.StatusBadRequest,
			wantBody:   "empty body\n",
		},
	}

	router := Router(h)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", strings.NewReader(tt.body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code, "handler вернул неверный статус код")
			assert.Equal(t, tt.wantBody, rr.Body.String(), "handler вернул неверное тело")
		})
	}
}

func TestGet(t *testing.T) {
	links := storage.NewStorage()
	config := config.NewConfig()
	h := &handlers.Handler{
		Storage: links,
		Config:  config,
	}

	h.Storage.Set("http://www.ya.ru/")

	tests := []struct {
		name         string
		url          string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "Valid ID",
			url:          "/6f782b56",
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "http://www.ya.ru/",
		},
		{
			name:         "Invalid ID",
			url:          "/123",
			wantStatus:   http.StatusNotFound,
			wantLocation: "",
		},
	}

	router := Router(h)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code, "handler вернул неверный код ответа")
			assert.Equal(t, tt.wantLocation, rr.Header().Get("Location"), "handler неверный локейшен")
		})
	}
}
