package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

// createTestServer creates a test server with multiple pages that have artificial delay
func createTestServer(delay time.Duration) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body>
			<h1>Home Page</h1>
			<p>Welcome to the home page.</p>
			<a href="/page1">Page 1</a>
			<a href="/page2">Page 2</a>
			<a href="/page3">Page 3</a>
			<a href="/page4">Page 4</a>
		</body></html>`))
	})

	mux.HandleFunc("/page1", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body>
			<h1>Page 1</h1>
			<p>This is page 1.</p>
			<a href="/">Home</a>
		</body></html>`))
	})

	mux.HandleFunc("/page2", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body>
			<h1>Page 2</h1>
			<p>This is page 2.</p>
			<a href="/">Home</a>
		</body></html>`))
	})

	mux.HandleFunc("/page3", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body>
			<h1>Page 3</h1>
			<p>This is page 3.</p>
			<a href="/">Home</a>
		</body></html>`))
	})

	mux.HandleFunc("/page4", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body>
			<h1>Page 4</h1>
			<p>This is page 4.</p>
			<a href="/">Home</a>
		</body></html>`))
	})

	return httptest.NewServer(mux)
}

func runCrawlWithConcurrency(serverURL string, maxConcurrency int) (map[string]PageData, time.Duration) {
	baseURL, _ := url.Parse(serverURL)

	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}

	start := time.Now()
	cfg.wg.Add(1)
	go cfg.crawlPage(serverURL)
	cfg.wg.Wait()
	elapsed := time.Since(start)

	return cfg.pages, elapsed
}

func TestCrawlConcurrency1(t *testing.T) {
	server := createTestServer(50 * time.Millisecond)
	defer server.Close()

	pages, elapsed := runCrawlWithConcurrency(server.URL, 1)

	if len(pages) != 5 {
		t.Errorf("expected 5 pages, got %d", len(pages))
	}

	t.Logf("Concurrency 1: crawled %d pages in %v", len(pages), elapsed)

	// With concurrency 1 and 50ms delay per page, 5 pages should take ~250ms minimum
	if elapsed < 200*time.Millisecond {
		t.Errorf("expected at least 200ms with concurrency 1, got %v", elapsed)
	}
}

func TestCrawlConcurrency2(t *testing.T) {
	server := createTestServer(50 * time.Millisecond)
	defer server.Close()

	pages, elapsed := runCrawlWithConcurrency(server.URL, 2)

	if len(pages) != 5 {
		t.Errorf("expected 5 pages, got %d", len(pages))
	}

	t.Logf("Concurrency 2: crawled %d pages in %v", len(pages), elapsed)
}

func TestCrawlConcurrency5(t *testing.T) {
	server := createTestServer(50 * time.Millisecond)
	defer server.Close()

	pages, elapsed := runCrawlWithConcurrency(server.URL, 5)

	if len(pages) != 5 {
		t.Errorf("expected 5 pages, got %d", len(pages))
	}

	t.Logf("Concurrency 5: crawled %d pages in %v", len(pages), elapsed)
}

func TestCrawlConcurrency10(t *testing.T) {
	server := createTestServer(50 * time.Millisecond)
	defer server.Close()

	pages, elapsed := runCrawlWithConcurrency(server.URL, 10)

	if len(pages) != 5 {
		t.Errorf("expected 5 pages, got %d", len(pages))
	}

	t.Logf("Concurrency 10: crawled %d pages in %v", len(pages), elapsed)
}

func TestCrawlDataConsistency(t *testing.T) {
	server := createTestServer(10 * time.Millisecond)
	defer server.Close()

	// Run with different concurrency levels and verify we get the same data
	pages1, _ := runCrawlWithConcurrency(server.URL, 1)
	pages2, _ := runCrawlWithConcurrency(server.URL, 2)
	pages5, _ := runCrawlWithConcurrency(server.URL, 5)
	pages10, _ := runCrawlWithConcurrency(server.URL, 10)

	// All should have the same number of pages
	if len(pages1) != len(pages2) || len(pages2) != len(pages5) || len(pages5) != len(pages10) {
		t.Errorf("page counts differ: 1=%d, 2=%d, 5=%d, 10=%d",
			len(pages1), len(pages2), len(pages5), len(pages10))
	}

	// All should have the same URLs
	for url := range pages1 {
		if _, exists := pages2[url]; !exists {
			t.Errorf("URL %s missing from concurrency 2 results", url)
		}
		if _, exists := pages5[url]; !exists {
			t.Errorf("URL %s missing from concurrency 5 results", url)
		}
		if _, exists := pages10[url]; !exists {
			t.Errorf("URL %s missing from concurrency 10 results", url)
		}
	}

	t.Logf("All concurrency levels returned consistent data with %d pages", len(pages1))
}

func TestCrawlPerformanceComparison(t *testing.T) {
	server := createTestServer(50 * time.Millisecond)
	defer server.Close()

	_, elapsed1 := runCrawlWithConcurrency(server.URL, 1)
	_, elapsed2 := runCrawlWithConcurrency(server.URL, 2)
	_, elapsed5 := runCrawlWithConcurrency(server.URL, 5)
	_, elapsed10 := runCrawlWithConcurrency(server.URL, 10)

	fmt.Printf("\n=== Concurrency Performance Comparison ===\n")
	fmt.Printf("Concurrency 1:  %v\n", elapsed1)
	fmt.Printf("Concurrency 2:  %v (%.2fx faster than 1)\n", elapsed2, float64(elapsed1)/float64(elapsed2))
	fmt.Printf("Concurrency 5:  %v (%.2fx faster than 1)\n", elapsed5, float64(elapsed1)/float64(elapsed5))
	fmt.Printf("Concurrency 10: %v (%.2fx faster than 1)\n", elapsed10, float64(elapsed1)/float64(elapsed10))

	// Concurrency 2 should be faster than concurrency 1
	if elapsed2 >= elapsed1 {
		t.Logf("Warning: Concurrency 2 (%v) not faster than concurrency 1 (%v)", elapsed2, elapsed1)
	}

	// Higher concurrency should generally be faster
	if elapsed5 >= elapsed2 {
		t.Logf("Warning: Concurrency 5 (%v) not faster than concurrency 2 (%v)", elapsed5, elapsed2)
	}
}
