package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type TestCase struct {
	name        string
	path        string
	event       string
	payloadFile string
	expCode     int
	expContent  string
}

func setupTestAPI(t *testing.T) (string, func()) {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(webhook))

	fmt.Printf("Running server on %s\n", ts.URL)

	return ts.URL, func() {
		ts.Close()
	}
}

func runTest(t *testing.T, tc TestCase, url string) {
	t.Run(tc.name, func(t *testing.T) {
		payloadFile, err := os.Open(tc.payloadFile)
		if err != nil {
			panic(err)
		}
		defer payloadFile.Close()

		req, err := http.NewRequest("POST", url+tc.path, payloadFile)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-GitHub-Event", tc.event)

		client := &http.Client{
			Timeout: time.Second * 10,
		}

		r, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer r.Body.Close()

		if r.StatusCode != tc.expCode {
			t.Fatalf("Expected %q, got %q.", http.StatusText(tc.expCode),
				http.StatusText(r.StatusCode))
		}
	})
}
