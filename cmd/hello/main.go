package main

import (
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

func main() {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		err = tmpl.Execute(w, data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
