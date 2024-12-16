package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// HTTP CONNECT tunnel establishment
func (p *Proxy) handleConnect(w http.ResponseWriter, r *http.Request) {
	log.Printf("Starting CONNECT tunnel: Host=%s, RemoteAddr=%s", r.Host, r.RemoteAddr)

	// Establish destination connection
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	destConn, err := (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext(ctx, "tcp", r.Host)
	if err != nil {
		log.Printf("Destination connection failed: %v", err)
		http.Error(w, fmt.Sprintf("Connection failed: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	// Hijack client connection
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Println("Hijacking not supported")
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Printf("Hijack Error: %v", err)
		http.Error(w, fmt.Sprintf("Hijack failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// Tunnel establishment response
	tunnelEstablished := "HTTP/1.1 200 Connection Established\r\n\r\n"
	_, err = clientConn.Write([]byte(tunnelEstablished))
	if err != nil {
		log.Printf("Tunnel Response Write Error: %v", err)
		return
	}

	// Bidirectional transfer
	transferCompleted := make(chan struct{}, 2)
	transferErrors := make(chan error, 2)

	transfer := func(dst io.Writer, src io.Reader, name string) {
		defer func() {
			transferCompleted <- struct{}{}
		}()

		written, err := io.Copy(dst, src)
		if err != nil && err != io.EOF {
			transferErrors <- fmt.Errorf("%s transfer error after %d bytes: %v", name, written, err)
			return
		}

		log.Printf("%s transfer completed successfully: %d bytes", name, written)
	}

	// Starting transfers
	go transfer(destConn, clientConn, "Client->Dest")
	go transfer(clientConn, destConn, "Dest->Client")

	// Waiting for transfers to complete or error out
	completedTransfers := 0
	var finalErr error

	for completedTransfers < 2 {
		select {
		case <-transferCompleted:
			completedTransfers++
		case err := <-transferErrors:
			finalErr = err
			log.Printf("Transfer error: %v", err)
		}
	}

	if finalErr != nil {
		log.Printf("Tunnel completed with error: %v", finalErr)
	} else {
		log.Println("Tunnel completed successfully")
	}
}
