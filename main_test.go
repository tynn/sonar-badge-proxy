package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServer_Addr(t *testing.T) {
	s := newServer()
	if s.Addr != ":4000" {
		t.Fatalf("wrong Addr=%s", s.Addr)
	}
}

func TestNewServer_Favicon(t *testing.T) {
	r, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := newServer().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	if code := w.Code; code != http.StatusNotFound {
		t.Errorf("wrong Code=%d for Path=%s", code, r.URL.Path)
	}
}
func TestNewServer_Root(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := newServer().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	if code := w.Code; code != http.StatusNotFound {
		t.Errorf("wrong Code=%d for Path=%s", code, r.URL.Path)
	}
}

func TestMain(t *testing.T) {
	go main()
	time.Sleep(time.Second)
}
