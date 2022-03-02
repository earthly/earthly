package main

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	go main()
	time.Sleep(time.Second) // Leave time for service to start
	resp, err := http.Get("http://localhost:8080/two/hello")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := "Hello, Friend!"
	actual, _ := ioutil.ReadAll(resp.Body)
	if expected != string(actual) {
		t.Fail()
	}
}
