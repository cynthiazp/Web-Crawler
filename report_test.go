package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteCSVReportBasic(t *testing.T) {
	pages := map[string]PageData{
		"example.com": {
			URL:            "https://example.com",
			H1:             "Example Title",
			FirstParagraph: "This is a paragraph.",
			OutgoingLinks:  []string{"https://example.com/page1", "https://example.com/page2"},
			ImageURLs:      []string{"https://example.com/img1.jpg"},
		},
	}

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_report.csv")

	err := writeCSVReport(pages, filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read and verify the CSV
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	// Check header
	expectedHeader := []string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"}
	if len(records) < 1 {
		t.Fatal("expected at least 1 row (header)")
	}
	for i, col := range expectedHeader {
		if records[0][i] != col {
			t.Errorf("header column %d: expected %q, got %q", i, col, records[0][i])
		}
	}

	// Check data row
	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + 1 data), got %d", len(records))
	}

	dataRow := records[1]
	if dataRow[0] != "example.com" {
		t.Errorf("expected page_url 'example.com', got %q", dataRow[0])
	}
	if dataRow[1] != "Example Title" {
		t.Errorf("expected h1 'Example Title', got %q", dataRow[1])
	}
	if dataRow[2] != "This is a paragraph." {
		t.Errorf("expected first_paragraph 'This is a paragraph.', got %q", dataRow[2])
	}
	if dataRow[3] != "https://example.com/page1;https://example.com/page2" {
		t.Errorf("expected outgoing_link_urls joined with ';', got %q", dataRow[3])
	}
	if dataRow[4] != "https://example.com/img1.jpg" {
		t.Errorf("expected image_urls 'https://example.com/img1.jpg', got %q", dataRow[4])
	}
}

func TestWriteCSVReportMultiplePages(t *testing.T) {
	pages := map[string]PageData{
		"example.com": {
			URL:            "https://example.com",
			H1:             "Home",
			FirstParagraph: "Welcome.",
			OutgoingLinks:  []string{"https://example.com/about"},
			ImageURLs:      []string{},
		},
		"example.com/about": {
			URL:            "https://example.com/about",
			H1:             "About",
			FirstParagraph: "About us.",
			OutgoingLinks:  []string{},
			ImageURLs:      []string{"https://example.com/logo.png", "https://example.com/team.jpg"},
		},
	}

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_report.csv")

	err := writeCSVReport(pages, filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	// Header + 2 data rows
	if len(records) != 3 {
		t.Errorf("expected 3 rows (header + 2 data), got %d", len(records))
	}
}

func TestWriteCSVReportEmptyPages(t *testing.T) {
	pages := map[string]PageData{}

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_report.csv")

	err := writeCSVReport(pages, filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	// Only header row
	if len(records) != 1 {
		t.Errorf("expected 1 row (header only), got %d", len(records))
	}
}

func TestWriteCSVReportEmptySlices(t *testing.T) {
	pages := map[string]PageData{
		"example.com": {
			URL:            "https://example.com",
			H1:             "Title",
			FirstParagraph: "Paragraph.",
			OutgoingLinks:  []string{},
			ImageURLs:      []string{},
		},
	}

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_report.csv")

	err := writeCSVReport(pages, filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	dataRow := records[1]
	// Empty slices should result in empty strings
	if dataRow[3] != "" {
		t.Errorf("expected empty outgoing_link_urls, got %q", dataRow[3])
	}
	if dataRow[4] != "" {
		t.Errorf("expected empty image_urls, got %q", dataRow[4])
	}
}

func TestWriteCSVReportInvalidPath(t *testing.T) {
	pages := map[string]PageData{}

	// Try to write to an invalid path
	err := writeCSVReport(pages, "/nonexistent/directory/report.csv")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
