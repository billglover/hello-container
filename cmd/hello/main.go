package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// pageInfo holds all the data used to render the index page.
type pageInfo struct {
	Version  string
	Commit   string
	Date     string
	Hostname string
}

// Server holds configuration details for the HTTP server.
type server struct {
	cancel func()
	tmpl   *template.Template
	s      http.Server
}

func main() {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := http.NewServeMux()
	s := server{cancel: cancel, tmpl: tmpl, s: http.Server{Addr: ":8080", Handler: m}}

	m.HandleFunc("/", s.info())
	m.HandleFunc("/kill", s.kill(s.info()))
	m.HandleFunc("/health", s.health())

	go func() {
		if err := s.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	select {
	case <-ctx.Done():
		s.s.Shutdown(ctx)
	}
	fmt.Fprintln(os.Stdout, "application terminated by user")
}

func (s *server) info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")

		hostname, err := os.Hostname()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		data := pageInfo{
			Version:  version,
			Commit:   commit,
			Date:     date,
			Hostname: hostname,
		}

		err = s.tmpl.Execute(w, data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func (s *server) health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (s *server) kill(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		defer s.cancel()
	}
}
