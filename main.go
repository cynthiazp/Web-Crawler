package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		fmt.Println("usage: crawler <url> [maxConcurrency] [maxPages]")
		os.Exit(1)
	}

	rawBaseURL := args[0]

	// Default values
	maxConcurrency := 10
	maxPages := 0 // 0 means no limit

	// Parse optional maxConcurrency
	if len(args) >= 2 {
		val, err := strconv.Atoi(args[1])
		if err != nil || val < 1 {
			fmt.Println("maxConcurrency must be a positive integer")
			os.Exit(1)
		}
		maxConcurrency = val
	}

	// Parse optional maxPages
	if len(args) >= 3 {
		val, err := strconv.Atoi(args[2])
		if err != nil || val < 0 {
			fmt.Println("maxPages must be a non-negative integer")
			os.Exit(1)
		}
		maxPages = val
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)
	fmt.Printf("  maxConcurrency: %d\n", maxConcurrency)
	fmt.Printf("  maxPages: %d (0 = unlimited)\n", maxPages)

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing base URL: %v\n", err)
		os.Exit(1)
	}

	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	fmt.Println("\n--- Crawl Results ---")
	fmt.Printf("Crawled %d pages\n", len(cfg.pages))

	// Write CSV report
	reportFile := "report.csv"
	if err := writeCSVReport(cfg.pages, reportFile); err != nil {
		fmt.Printf("error writing CSV report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Report written to: %s\n", reportFile)
}