package main

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func AssertDeepEqual(t *testing.T, e interface{}, a interface{}, m string) {
	if !reflect.DeepEqual(e, a) {
		t.Errorf(m, a)
	}
}

func AssertEqual(t *testing.T, e interface{}, a interface{}, m string) {
	if a != e {
		t.Errorf(m, a)
	}
}

func AssertNil(t *testing.T, a interface{}, m string) {
	if a != nil {
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

func AssertPointerEqual(t *testing.T, e interface{}, a interface{}, m string) {
	if reflect.ValueOf(a).Pointer() != reflect.ValueOf(e).Pointer() {
		t.Errorf(m, a)
	}
}

func TestMain(t *testing.T) {
	os.Setenv("PORT", "4000")
	os.Setenv("REMOTE", "sonarcloud.io")
	os.Setenv("SECRET", "012345789abcdef")
	os.Setenv("METRIC", "bugs,lines")
	go main()
	time.Sleep(time.Second)
}
