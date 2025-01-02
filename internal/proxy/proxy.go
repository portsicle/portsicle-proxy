package proxy

import (
	"log"
	"net/http"
	"strings"

	"github.com/amitsuthar69/portsicle-proxy/internal/database"
)

// proxy configuration
type Config struct {
	ListenAddr string
}

// proxy server
type Proxy struct {
	config     Config
	httpServer *http.Server
}

// NewProxy will create a new proxy instance
func NewProxy(cfg Config) *Proxy {
	return &Proxy{
		config: cfg,
	}
}

func (p *Proxy) Start() error {
	server := &http.Server{
		Addr: p.config.ListenAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// header obfuscation before processing the request
			obfuscateHeaders(r)

			log.Printf("Incoming request: Method=%s, Host=%s, URL=%s",
				r.Method, r.Host, r.URL)

			db, err := database.InitDB()
			if err != nil {
				log.Fatalf("Error connecting to database %v", err)
			}

			blocked_domains, err := database.GetBlockedDomains(db)
			if err != nil {
				log.Fatalf("Error getting blocked domains %v", err)
			}

			// Split r.Host to get the domain without port
			hostParts := strings.Split(r.Host, ":")
			normalizedHost := strings.ToLower(hostParts[0])

			// check if the host is a user blocked domain
			if isBlockedDomain(blocked_domains, normalizedHost) {
				log.Printf("domain %v is blocked", r.Host)
				http.Error(w, "Access restricted to this Domain!", http.StatusForbidden)
				return
			}

			if r.Method == "CONNECT" {
				p.handleConnect(w, r)
				return
			}

			log.Printf("Non-CONNECT method received: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}),
	}

	p.httpServer = server
	return server.ListenAndServe()
}

func (p *Proxy) Stop() error {
	if p.httpServer != nil {
		return p.httpServer.Close()
	}
	return nil
}

func isBlockedDomain(blockedDomains map[string]struct{}, normalizedHost string) bool {
	// Direct domain match
	if _, exists := blockedDomains[normalizedHost]; exists {
		return true
	}

	// Check parent domains within sub domain parts
	parts := strings.Split(normalizedHost, ".")
	for i := 1; i < len(parts); i++ {
		parentDomain := strings.Join(parts[i:], ".")
		if _, exists := blockedDomains[parentDomain]; exists {
			return true
		}
	}

	return false
}
