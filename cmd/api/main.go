package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
	})
	mux.HandleFunc("/records/gender", func(w http.ResponseWriter, r *http.Request) {
	})
	mux.HandleFunc("/records/birthdate", func(w http.ResponseWriter, r *http.Request) {
	})
	mux.HandleFunc("/records/name", func(w http.ResponseWriter, r *http.Request) {
	})
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
}
