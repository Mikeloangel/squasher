package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/state"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/config"
	"github.com/Mikeloangel/squasher/internal/models"
	"github.com/stretchr/testify/assert"
)

// Tests handlers for creation of short url
func TestCreateShortURL(t *testing.T) {
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

	h := getHandlers()
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

// Tests get original url
func TestGetOriginalURL(t *testing.T) {
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

	h := getHandlers()
	h.Storage.StoreURL("http://www.ya.ru/")
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

// Tests handlers for creation of short url via api and json format
func TestHandlerCreateShortURLJson(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CreateShortURLRequest
		wantStatus   int
		wantResponse models.CreateShortURLResponse
	}{
		{
			name: "Valid URL",
			body: models.CreateShortURLRequest{
				URL: "http://www.ya.ru/",
			},
			wantStatus: http.StatusCreated,
			wantResponse: models.CreateShortURLResponse{
				Result: "http://localhost:8080/6f782b56",
			},
		},
	}

	h := getHandlers()
	router := Router(h)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := json.Marshal(tt.body)
			req, err := http.NewRequest("POST", "/api/shorten", strings.NewReader(string(request)))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code, "hanler вернул неверный статус код")
			s := models.CreateShortURLResponse{}

			json.Unmarshal(rr.Body.Bytes(), &s)
			assert.Equal(t, tt.wantResponse.Result, s.Result, "handler вернул неверное тело")
		})
	}
}

// Sets up handlers with app state and configuration
func getHandlers() *handlers.Handler {
	var err error

	cfg := config.NewConfig(
		"localhost",
		8080,
		"http://localhost:8080",
		"/tmp/short-url-db.json",
		"",
	)

	storage := storage.NewInMemoryStorage()
	err = storage.Init()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return handlers.NewHandler(
		state.NewState(storage, cfg, &sql.DB{}),
	)
}
