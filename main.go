package main

import "log"

func logFatalPanic() { log.Fatal(recover()) }

func main() {
	defer logFatalPanic()
	c := LoadConfig()
	p := NewProxy(c)
	s := p.Server()
	err := s.ListenAndServe()
	panic(err)
}
