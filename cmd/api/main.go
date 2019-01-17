package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/lag13/records/internal/db"
	"github.com/lag13/records/internal/endpoints/getsortperson"
	"github.com/lag13/records/internal/endpoints/postrecord"
	"github.com/lag13/records/internal/person"
)

var mu = &sync.Mutex{}

func writeAndLogErr(w http.ResponseWriter, body []byte) {
	if _, err := w.Write(body); err != nil {
		log.Print(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
		p, resp, err := postrecord.PostRecord(r)
		if err != nil {
			log.Print(err)
		}
		w.WriteHeader(resp.StatusCode)
		if len(resp.Errors) == 0 {
			mu.Lock()
			defer mu.Unlock()
			db.Persons = append(db.Persons, p)
			return
		}
		// TODO: This json encoding logic is duplicated in
		// other places and should be consolidated. More
		// generally, the logic of writing the response is
		// duplicated and could be consolidated.
		body, err := json.Marshal(resp)
		if err != nil {
			// The only time json.Marshal fails is if a
			// type is passed in which cannot be
			// marshalled so to me a panic is acceptable
			// here.
			panic(err)
		}
		writeAndLogErr(w, body)
	})
	mux.HandleFunc("/records/gender", func(w http.ResponseWriter, r *http.Request) {
		resp := getsortperson.Sort(r, person.SortGenderLastNameAsc, db.Persons)
		w.WriteHeader(resp.StatusCode)
		body, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		writeAndLogErr(w, body)
	})
	mux.HandleFunc("/records/birthdate", func(w http.ResponseWriter, r *http.Request) {
		resp := getsortperson.Sort(r, person.SortBirthdateAsc, db.Persons)
		w.WriteHeader(resp.StatusCode)
		body, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		writeAndLogErr(w, body)
	})
	mux.HandleFunc("/records/name", func(w http.ResponseWriter, r *http.Request) {
		resp := getsortperson.Sort(r, person.SortLastNameDesc, db.Persons)
		w.WriteHeader(resp.StatusCode)
		body, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		writeAndLogErr(w, body)
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
