package main

import (
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

// addPageVisit checks if a page has been visited and adds it if not
func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, exists := cfg.pages[normalizedURL]; exists {
		return false
	}

	cfg.pages[normalizedURL] = PageData{}
	return true
}

// pagesLen returns the current number of pages (thread-safe)
func (cfg *config) pagesLen() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages)
}

// crawlPage recursively crawls pages starting from rawCurrentURL
func (cfg *config) crawlPage(rawCurrentURL string) {
	// Acquire a slot from the concurrency control channel
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	// Check if we've reached the max pages limit
	if cfg.maxPages > 0 && cfg.pagesLen() >= cfg.maxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	// Only crawl pages on the same domain
	if cfg.baseURL.Host != currentURL.Host {
		return
	}

	// Normalize the current URL
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	// If we've already visited this page, return
	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst {
		return
	}

	// Fetch the HTML
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	// Extract page data and update the map
	pageData := extractPageData(html, rawCurrentURL)
	cfg.mu.Lock()
	cfg.pages[normalizedURL] = pageData
	cfg.mu.Unlock()

	// Crawl each link concurrently
	for _, link := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(link)
	}
}
