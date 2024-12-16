package proxy

import (
	"net/http"
	"strings"
)

// sanitizeHost will remove port information from host
func sanitizeHost(host string) string {
	parts := strings.Split(host, ":")
	return parts[0]
}

// obfuscateHeaders will remove or modify potentially identifying headers
func obfuscateHeaders(r *http.Request) {
	// Removing potentially identifying headers
	headersToRemove := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"Forwarded",
		"X-Forwarded-Host",
		"X-Forwarded-Proto",
		"Client-IP",
		"X-Client-IP",
		"Referer",
	}

	for _, header := range headersToRemove {
		r.Header.Del(header)
	}

	// additional privacy headers
	r.Header.Set("DNT", "1") // Do Not Track
	r.Header.Set("Sec-Fetch-Dest", "document")
	r.Header.Set("Sec-Fetch-Mode", "navigate")
	r.Header.Set("Sec-Fetch-Site", "none")
	r.Header.Set("Sec-Fetch-User", "?1")

	// Standardizing User-Agent
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")
}
