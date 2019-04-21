package main

import (
	"log"
	"net/http"
)

func logFatalPanic() { log.Fatal(recover()) }

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}
func newServer() *http.Server {
	c := LoadConfig()
	m := http.NewServeMux()
	m.HandleFunc("/favicon.ico", faviconHandler)
	m.HandleFunc("/", http.NotFound)
	return &http.Server{Addr: c.Addr, Handler: m}
}

func main() {
	defer logFatalPanic()
	panic(newServer().ListenAndServe())
}
