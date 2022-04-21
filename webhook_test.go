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

func TestPost(t *testing.T) {
	testCases := []TestCase{
		{name: "01 workflow run requested", path: "/", event: "workflow_run", payloadFile: "01_workflow_run_2096567160_requested.json", expCode: http.StatusOK, expContent: ""},
		{name: "02 workflow job queued", path: "/", event: "workflow_job", payloadFile: "02_workflow_job_5834689160_queued.json", expCode: http.StatusOK, expContent: ""},
		{name: "03 check run created", path: "/", event: "check_run", payloadFile: "03_check_run_5834689160_created.json", expCode: http.StatusOK, expContent: ""},
		{name: "04 workflow job completed", path: "/", event: "workflow_job", payloadFile: "04_workflow_job_5834689160_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "05 check run completed", path: "/", event: "check_run", payloadFile: "05_check_run_5834689160_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "06 workflow job queued", path: "/", event: "workflow_job", payloadFile: "06_workflow_job_5834695996_queued.json", expCode: http.StatusOK, expContent: ""},
		{name: "07 check run created", path: "/", event: "check_run", payloadFile: "07_check_run_5834695996_created.json", expCode: http.StatusOK, expContent: ""},
		{name: "08 deployment created", path: "/", event: "deployment", payloadFile: "08_deployment_540578780_created.json", expCode: http.StatusOK, expContent: ""},
		{name: "09 workflow job queued", path: "/", event: "workflow_job", payloadFile: "09_workflow_job_5834696229_queued.json", expCode: http.StatusOK, expContent: ""},
		{name: "10 check run created", path: "/", event: "check_run", payloadFile: "10_check_run_5834696229_created.json", expCode: http.StatusOK, expContent: ""},
		{name: "11 workflow job in progress", path: "/", event: "workflow_job", payloadFile: "11_workflow_job_5834695996_in_progress.json", expCode: http.StatusOK, expContent: ""},
		{name: "12 workflow job in progress", path: "/", event: "workflow_job", payloadFile: "12_workflow_job_5834696229_in_progress.json", expCode: http.StatusOK, expContent: ""},
		{name: "13 deployment status created", path: "/", event: "deployment_status", payloadFile: "13_deployment_status_540578780_created.json", expCode: http.StatusOK, expContent: ""},
		{name: "14 check run completed", path: "/", event: "check_run", payloadFile: "14_check_run_5834695996_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "15 workflow job completed", path: "/", event: "workflow_job", payloadFile: "15_workflow_job_5834695996_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "16 workflow job completed", path: "/", event: "workflow_job", payloadFile: "16_workflow_job_5834696229_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "17 check run completed", path: "/", event: "check_run", payloadFile: "17_check_run_5834696229_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "18 deployment status created", path: "/", event: "deployment_status", payloadFile: "18_deployment_status_540578780_created.json", expCode: http.StatusOK, expContent: ""},
		{name: "19 check suite completed", path: "/", event: "check_suite", payloadFile: "19_check_suite_5939269259_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "20 workflow run completed", path: "/", event: "workflow_run", payloadFile: "20_workflow_run_2096567160_completed.json", expCode: http.StatusOK, expContent: ""},
	}

	initialize()

	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, tc := range testCases {
		runTest(t, tc, url)
	}
}

func runTest(t *testing.T, tc TestCase, url string) {
	t.Run(tc.name, func(t *testing.T) {
		payloadFile, err := os.Open("test/" + tc.payloadFile)
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

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(webhook))

	fmt.Printf("Running server on %s\n", ts.URL)

	return ts.URL, func() {
		ts.Close()
	}
}
