package main

import (
	"log"
	"net/http"
)

func logFatalPanic() { log.Fatal(recover()) }

func newServer() *http.Server {
	c := LoadConfig()
	m := http.NewServeMux()
	m.HandleFunc("/favicon.ico", http.NotFound)
	m.HandleFunc("/", http.NotFound)
	return &http.Server{Addr: c.Addr, Handler: m}
}

func main() {
	defer logFatalPanic()
	panic(newServer().ListenAndServe())
}
