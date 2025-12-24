package main

import (
	"encoding/csv"
	"os"
	"strings"
)

// writeCSVReport writes the crawled pages data to a CSV file
func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for pageURL, pageData := range pages {
		row := []string{
			pageURL,
			pageData.H1,
			pageData.FirstParagraph,
			strings.Join(pageData.OutgoingLinks, ";"),
			strings.Join(pageData.ImageURLs, ";"),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
