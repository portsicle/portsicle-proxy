package main

import (
	"log"
	"os"

	"github.com/amitsuthar69/attorney-toolkit/internal/proxy"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Proxy with default config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	p := proxy.NewProxy(proxy.Config{
		ListenAddr: ":" + port,
	})

	if err := p.Start(); err != nil {
		log.Fatalf("Proxy startup failed: %v", err)
	}
}
