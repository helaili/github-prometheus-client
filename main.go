package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v43/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const github_actions_workflow_run_total_str string = "github_actions_workflow_run_total"

var counters = make(map[string]*prometheus.CounterVec)

func main() {
	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	counters[github_actions_workflow_run_total_str] = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "github_actions_workflow_run_total",
		Help: "The total number of workflow runs",
	},
		[]string{"org", "repo", "workflow"},
	)

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
func webhook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s event on end point %s\n", r.Header.Get("X-GitHub-Event"), r.URL)

	// TODO: support webhook secret
	payload, err := github.ValidatePayload(r, []byte(""))
	if err != nil {
		log.Printf("error reading request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.WorkflowRunEvent:
		reportWorkflowRunEvent(e)
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func reportWorkflowRunEvent(event *github.WorkflowRunEvent) {
	log.Printf("reporting workflow run event for %s/%s/%s\n", event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName())
	github_actions_workflow_run_total, found := counters[github_actions_workflow_run_total_str]

	if !found {
		log.Printf("metric not registered %s\n", github_actions_workflow_run_total_str)
		return
	}
	github_actions_workflow_run_total.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName()).Inc()
}
