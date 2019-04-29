package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/enter", enterHandler)
	http.HandleFunc("/finish", finishHandler)

	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func enterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/enter" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintln(w, "EnterHandler!")
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/finish" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintln(w, "FinishHandler!")
}
