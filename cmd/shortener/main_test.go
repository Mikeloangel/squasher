package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
)

func TestPost(t *testing.T) {
	links := storage.NewStorage()

	h := &handlers.Handler{
		Storage: links,
	}

	body := strings.NewReader("http://www.ya.ru/")
	rPost, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Post)
	handler.ServeHTTP(w, rPost)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	wantBody := "http://localhost:8080/6f782b56"
	if body := w.Body.String(); body != wantBody {
		t.Errorf("handler returned wrong body: got %s want %s", body, wantBody)
	}
}

func TestGet(t *testing.T) {
	links := storage.NewStorage()

	h := &handlers.Handler{
		Storage: links,
	}

	h.Storage.Set("http://www.ya.ru/")

	rPost, err := http.NewRequest("GET", "/6f782b56", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Get)
	handler.ServeHTTP(w, rPost)

	if status := w.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	wantLocation := "http://www.ya.ru/"
	if location := w.Header().Get("Location"); location != wantLocation {
		t.Errorf("handler returned wrong location: got %s want %s", location, wantLocation)
	}
}
