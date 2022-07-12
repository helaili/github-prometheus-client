package main

import (
	"net/http"
	"testing"
)

func TestWorkflow(t *testing.T) {
	testCases := []TestCase{
		{name: "01 workflow run requested", path: "/", event: "workflow_run", payloadFile: "test/workflow/01_workflow_run_2096567160_requested.json", expCode: http.StatusOK, expContent: ""},
		{name: "02 workflow job queued", path: "/", event: "workflow_job", payloadFile: "test/workflow/02_workflow_job_5834689160_queued.json", expCode: http.StatusOK, expContent: ""},
		{name: "03 check run created", path: "/", event: "check_run", payloadFile: "test/workflow/03_check_run_5834689160_created.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "04 workflow job completed", path: "/", event: "workflow_job", payloadFile: "test/workflow/04_workflow_job_5834689160_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "05 check run completed", path: "/", event: "check_run", payloadFile: "test/workflow/05_check_run_5834689160_completed.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "06 workflow job queued", path: "/", event: "workflow_job", payloadFile: "test/workflow/06_workflow_job_5834695996_queued.json", expCode: http.StatusOK, expContent: ""},
		{name: "07 check run created", path: "/", event: "check_run", payloadFile: "test/workflow/07_check_run_5834695996_created.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "08 deployment created", path: "/", event: "deployment", payloadFile: "test/workflow/08_deployment_540578780_created.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "09 workflow job queued", path: "/", event: "workflow_job", payloadFile: "test/workflow/09_workflow_job_5834696229_queued.json", expCode: http.StatusOK, expContent: ""},
		{name: "10 check run created", path: "/", event: "check_run", payloadFile: "test/workflow/10_check_run_5834696229_created.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "11 workflow job in progress", path: "/", event: "workflow_job", payloadFile: "test/workflow/11_workflow_job_5834695996_in_progress.json", expCode: http.StatusOK, expContent: ""},
		{name: "12 workflow job in progress", path: "/", event: "workflow_job", payloadFile: "test/workflow/12_workflow_job_5834696229_in_progress.json", expCode: http.StatusOK, expContent: ""},
		{name: "13 deployment status created", path: "/", event: "deployment_status", payloadFile: "test/workflow/13_deployment_status_540578780_created.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "14 check run completed", path: "/", event: "check_run", payloadFile: "test/workflow/14_check_run_5834695996_completed.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "15 workflow job completed", path: "/", event: "workflow_job", payloadFile: "test/workflow/15_workflow_job_5834695996_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "16 workflow job completed", path: "/", event: "workflow_job", payloadFile: "test/workflow/16_workflow_job_5834696229_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "17 check run completed", path: "/", event: "check_run", payloadFile: "test/workflow/17_check_run_5834696229_completed.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "18 deployment status created", path: "/", event: "deployment_status", payloadFile: "test/workflow/18_deployment_status_540578780_created.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "19 check suite completed", path: "/", event: "check_suite", payloadFile: "test/workflow/19_check_suite_5939269259_completed.json", expCode: http.StatusNotImplemented, expContent: ""},
		{name: "20 workflow run completed", path: "/", event: "workflow_run", payloadFile: "test/workflow/20_workflow_run_2096567160_completed.json", expCode: http.StatusOK, expContent: ""},
		{name: "01  installation created", path: "/", event: "installation", payloadFile: "test/installation/01_installation_25140335_created.json", expCode: http.StatusOK, expContent: ""},
	}

	_, _, _, app_id := initializeEnv()
	installationHandlers = make(map[string]*InstallationHandler)
	if cache == nil {
		cache = NewOfflineCache(app_id)
	}
	initializeDummyInstallationHandlers()

	url, cleanup := setupTestAPI(t)
	defer cleanup()

	for _, tc := range testCases {
		runTest(t, tc, url)
	}
}
