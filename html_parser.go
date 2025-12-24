 package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// PageData represents extracted data from a web page
type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

// getH1FromHTML extracts the text content of the first <h1> tag from HTML
func getH1FromHTML(htmlBody string) string {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return ""
	}

	var h1Text string
	var findH1 func(*html.Node)
	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			h1Text = extractText(n)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if h1Text != "" {
				return
			}
			findH1(c)
		}
	}
	findH1(doc)
	return h1Text
}

// getFirstParagraphFromHTML extracts the first <p> tag text, prioritizing <main> content
func getFirstParagraphFromHTML(htmlBody string) string {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return ""
	}

	// First, try to find <p> inside <main>
	mainNode := findNode(doc, "main")
	if mainNode != nil {
		pNode := findNode(mainNode, "p")
		if pNode != nil {
			return strings.TrimSpace(extractText(pNode))
		}
	}

	// Fallback: find first <p> anywhere
	pNode := findNode(doc, "p")
	if pNode != nil {
		return strings.TrimSpace(extractText(pNode))
	}

	return ""
}

// findNode finds the first node with the given tag name
func findNode(n *html.Node, tag string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findNode(c, tag); result != nil {
			return result
		}
	}
	return nil
}

// extractText extracts all text content from a node
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var result strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result.WriteString(extractText(c))
	}
	return result.String()
}

// getURLsFromHTML extracts all URLs from anchor tags in the HTML
func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		parsedURL, err := url.Parse(href)
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(parsedURL)
		urls = append(urls, resolvedURL.String())
	})

	return urls, nil
}

// getImagesFromHTML extracts all image URLs from img tags in the HTML
func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var images []string
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists || src == "" {
			return
		}

		parsedURL, err := url.Parse(src)
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(parsedURL)
		images = append(images, resolvedURL.String())
	})

	return images, nil
}

// extractPageData extracts all relevant data from a web page
func extractPageData(htmlBody string, rawURL string) PageData {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return PageData{URL: rawURL}
	}

	outgoingLinks, _ := getURLsFromHTML(htmlBody, baseURL)
	imageURLs, _ := getImagesFromHTML(htmlBody, baseURL)

	return PageData{
		URL:            rawURL,
		H1:             getH1FromHTML(htmlBody),
		FirstParagraph: getFirstParagraphFromHTML(htmlBody),
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
