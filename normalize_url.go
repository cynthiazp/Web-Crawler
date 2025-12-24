package main

import (
	"net/url"
	"strings"
)

// normalizeURL removes the scheme, query parameters, and fragments from a URL, but retains "www." if present.
func normalizeURL(input string) (string, error) {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	host := parsedURL.Host
	path := strings.TrimSuffix(parsedURL.Path, "/")

	return host + path, nil
}