package main

import (
	"net/http"
	"testing"
)

func TestInstallation(t *testing.T) {
	testCases := []TestCase{
		{name: "01  installation created", path: "/", event: "installation", payloadFile: "test/installation/01_installation_25140335_created.json", expCode: http.StatusOK, expContent: ""},
	}

	initializeEnv()
	initializeInstallationHandler()

	url, cleanup := setupTestAPI(t)
	defer cleanup()

	for _, tc := range testCases {
		runTest(t, tc, url)
	}
}
