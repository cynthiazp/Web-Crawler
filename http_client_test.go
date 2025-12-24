package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHTMLSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body><h1>Test</h1></body></html>"))
	}))
	defer server.Close()

	html, err := getHTML(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, "<h1>Test</h1>") {
		t.Errorf("expected HTML to contain <h1>Test</h1>, got %s", html)
	}
}

func TestGetHTMLErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := getHTML(server.URL)
	if err == nil {
		t.Fatal("expected error for 404 status, got nil")
	}

	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected error to contain 404, got %v", err)
	}
}

func TestGetHTMLWrongContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "hello"}`))
	}))
	defer server.Close()

	_, err := getHTML(server.URL)
	if err == nil {
		t.Fatal("expected error for wrong content type, got nil")
	}

	if !strings.Contains(err.Error(), "content type") {
		t.Errorf("expected error to mention content type, got %v", err)
	}
}

func TestGetHTMLUserAgent(t *testing.T) {
	var receivedUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUserAgent = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html></html>"))
	}))
	defer server.Close()

	_, err := getHTML(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedUserAgent != "BootCrawler/1.0" {
		t.Errorf("expected User-Agent 'BootCrawler/1.0', got %s", receivedUserAgent)
	}
}

func TestGetHTMLInvalidURL(t *testing.T) {
	_, err := getHTML("not-a-valid-url")
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}
