package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/lag13/records/internal/multicsv"
	"github.com/lag13/records/internal/person"
)

var db = []person.Person{}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
		// TODO: I wrote this code quickly to get a sense of
		// what needs to happen and I plan to reorganize and
		// unit test it. My problem though is that I don't
		// like it when a units from the same repository
		// reference eachother because if one unit breaks then
		// the other will too. So, I'm not sure how to
		// structure this code into a unit because it
		// definitely needs to reference other units in this
		// repository. Perhaps I make the unit accept two
		// functions which encapsulate what we need from
		// multicsv and person then in the real scenario we
		// pass them in, but I'm not sure if I like that extra
		// layer of abstraction. Or maybe we use the functions
		// directly but only test for the presence of errors
		// and not the exact wording (we will still test that
		// the happy path transformation works as expected but
		// I'm more okay with that since it probably won't
		// change much if at all). Or maybe this stuff is so
		// simple that covering it with e2e tests is
		// sufficient?
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`nopity nope`))
			return
		}
		lines, parseErrs := multicsv.ReadAll("|, ", 5, r.Body)
		if len(parseErrs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("got some errors %s", parseErrs[0])))
			return
		}
		p, parseErrs := person.Parse(lines[0])
		if len(parseErrs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("got some errors %s", strings.Join(parseErrs, ", "))))
			return
		}
		db = append(db, p)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/records/gender", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Probably should copy the db here so this GET
		// isn't modifying anything.
		person.SortGenderLastNameAsc(db)
		b, err := json.Marshal(db)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
	mux.HandleFunc("/records/birthdate", func(w http.ResponseWriter, r *http.Request) {
		person.SortBirthdateAsc(db)
		b, err := json.Marshal(db)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
	mux.HandleFunc("/records/name", func(w http.ResponseWriter, r *http.Request) {
		person.SortLastNameDesc(db)
		b, err := json.Marshal(db)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
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
