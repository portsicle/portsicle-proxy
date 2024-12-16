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
	r.Header.Del("X-Forwarded-For")
	r.Header.Del("X-Real-IP")

	// Standardizing User-Agent
	r.Header.Set("User-Agent", "Mozilla/5.0 (Anonymous Proxy)")
}
