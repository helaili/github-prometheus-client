package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type InstallationHandler struct {
	workflowRunMetrics *WorkflowRunMetrics
	workflowJobMetrics *WorkflowJobMetrics
}

func NewInstallationHandler(installation_id int64, cache IWorkflowNameCache) *InstallationHandler {
	inst := new(InstallationHandler)
	registry := prometheus.NewRegistry()
	inst.workflowRunMetrics = NewWorkflowRunMetrics(registry, cache)
	inst.workflowJobMetrics = NewWorkflowJobMetrics(registry, cache)

	handlerPath := fmt.Sprintf("/metrics-%d", installation_id)

	log.Printf("Receiving events for installation %d on path %s\n", installation_id, handlerPath)
	// This is the Prometheus endpoint.
	http.Handle(handlerPath, promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	return inst
}
