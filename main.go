package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v43/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var workflowNames = NewWorkflowNameCacheImpl()
var workflowRunMetrics = NewWorkflowRunMetrics(workflowNames)
var workflowJobMetrics = NewWorkflowJobMetrics(workflowNames)

func main() {
	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	// This is the Prometheus endpoint.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	// This is the GitHub Webhook endpoint.
	http.HandleFunc("/webhook", webhook)

	fmt.Printf("Listening on address %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

/*
 * Receive a webhook from GitHub and report the event to Prometheus.
 */
func webhook(w http.ResponseWriter, req *http.Request) {
	log.Printf("Received %s event on end point %s\n", req.Header.Get("X-GitHub-Event"), req.URL)

	// TODO: support webhook secret
	payload, err := github.ValidatePayload(req, []byte(""))
	if err != nil {
		log.Printf("error reading request body: err=%s\n", err)
		return
	}
	defer req.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(req), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.WorkflowRunEvent:
		workflowRunMetrics.report(github.WebHookType(req), e)
	case *github.WorkflowJobEvent:
		workflowJobMetrics.report(github.WebHookType(req), e)
	default:
		// log.Printf("unknown event type %s\n", github.WebHookType(req))
		return
	}

	w.WriteHeader(http.StatusOK)
}
