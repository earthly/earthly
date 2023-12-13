package main

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	go main()
	time.Sleep(time.Second) // Leave time for service to start
	resp, err := http.Get("http://localhost:8080/one/hello")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := "Hello, World!"
	actual, _ := io.ReadAll(resp.Body)
	if expected != string(actual) {
		t.Fail()
	}
}
