package proxy

import (
	"log"
	"net/http"
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
			log.Printf("Incoming request: Method=%s, Host=%s, URL=%s",
				r.Method, r.Host, r.URL)

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
