package proxy

import (
	"bufio"
	"bytes"
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

		reader := bufio.NewReader(src)
		buffer := make([]byte, 32*1024)

		for {
			n, err := reader.Read(buffer)
			if err != nil {
				if err != io.EOF {
					transferErrors <- fmt.Errorf("%s transfer error: %v", name, err)
				}
				return
			}

			data := buffer[:n]
			if name == "Dest->Client" {

				if bytes.Contains(data, []byte("HTTP/1.")) &&
					bytes.Contains(data, []byte("Content-Type: text/html")) {

					/*
						we will never reach here!
						the data is scrambled. I tried utf8 decoding and gzip decompression, but none worked!
						I think as the 'data' inside this CONNECT tunnel is TLS encrypted, above string comparion fails.
					*/

					parts := bytes.Split(data, []byte("\r\n\r\n"))
					if len(parts) > 1 {
						headers := parts[0]
						body := bytes.Join(parts[1:], []byte("\r\n\r\n"))
						if newBody, err := removeAds(body); err == nil {

							log.Print("new: ", newBody)

							headerLines := bytes.Split(headers, []byte("\r\n"))
							for i, line := range headerLines {
								if bytes.HasPrefix(bytes.ToLower(line), []byte("content-length:")) {
									headerLines[i] = []byte(fmt.Sprintf("Content-Length: %d", len(newBody)))
									break
								}
							}
							headers = bytes.Join(headerLines, []byte("\r\n"))

							data = append(headers, []byte("\r\n\r\n")...)
							data = append(data, newBody...)
						}
					}
				}
			}

			_, err = dst.Write(data)
			// log.Printf("%s transfer completed successfully: %d bytes", name, written)
			if err != nil {
				transferErrors <- fmt.Errorf("%s write error: %v", name, err)
				return
			}
		}
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
