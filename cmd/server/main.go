package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func delayedReply(delay time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-timer.C:
			fmt.Fprintf(w, "Ответ через %s\n", delay)
		case <-r.Context().Done():
			return
		}
	}
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "адрес сервера")
	flag.Parse()
	server := &http.Server{
		Addr:              *addr,
		Handler:           routes(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("Anime ASCII Club: http://%s/art — Ctrl+C to stop", *addr)
	log.Fatal(server.ListenAndServe())
}
