package main

import (
	"log"
	"net/http"
)

func logFatalPanic() { log.Fatal(recover()) }

func newServer() *http.Server {
	m := http.NewServeMux()
	m.HandleFunc("/favicon.ico", http.NotFound)
	m.HandleFunc("/", http.NotFound)
	return &http.Server{Addr: ":4000", Handler: m}
}

func main() {
	defer logFatalPanic()
	panic(newServer().ListenAndServe())
}
