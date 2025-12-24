package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "host name only",
			inputURL: "https://www.wikipedia.org",
			expected: "www.wikipedia.org",
		},
		{
			name:     "remove trailing slash",
			inputURL: "https://www.google.com/",
			expected: "www.google.com",
		},
		{
			name:     "handle subdomains",
			inputURL: "https://maps.google.com/path",
			expected: "maps.google.com/path",
		},
		{
			name:     "remove query parameters",
			inputURL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expected: "www.youtube.com/watch",
		},
		{
			name:     "remove fragment",
			inputURL: "https://en.wikipedia.org/wiki/URL#Syntax",
			expected: "en.wikipedia.org/wiki/URL",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}