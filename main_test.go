package main

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func AssertEqual(t *testing.T, e interface{}, a interface{}, m string) {
	if a != e {
		t.Errorf(m, a)
	}
}

func AssertDeepEqual(t *testing.T, e interface{}, a interface{}, m string) {
	if !reflect.DeepEqual(e, a) {
		t.Errorf(m, a)
	}
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func AssertPanic(t *testing.T, m string) {
	if err := recover(); err == nil {
		t.Errorf("No panic on %v", m)
	}
}

func ReadDotenv() {
	f, err := os.Open(".env")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := regexp.MustCompile("^export ([A-Z]+)=\"(.*)\"$")
	s := bufio.NewScanner(f)
	for s.Scan() {
		env := r.FindStringSubmatch(s.Text())
		if env != nil {
			os.Setenv(env[1], env[2])
		}
	}
}

func TestNewServer_Addr(t *testing.T) {
	ReadDotenv()
	s := newServer()
	AssertEqual(t, ":4000", s.Addr, "Wrong Addr=%v")
}

func TestNewServer_Favicon(t *testing.T) {
	ReadDotenv()
	r, err := http.NewRequest("GET", "/favicon.ico", nil)
	AssertNoError(t, err)

	h := newServer().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusNotFound, w.Code, "wrong Code=%v")
}

func TestNewServer_Root(t *testing.T) {
	ReadDotenv()
	r, err := http.NewRequest("GET", "/", nil)
	AssertNoError(t, err)

	h := newServer().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusNotFound, w.Code, "wrong Code=%v")
}

func TestMain(t *testing.T) {
	ReadDotenv()
	go main()
	time.Sleep(time.Second)
}
