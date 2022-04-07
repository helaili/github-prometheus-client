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

var counters = make(map[string]*prometheus.CounterVec)
var histograms = make(map[string]*prometheus.HistogramVec)
var workflowNames = make(map[string]string)

func main() {
	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	counters["github_actions_workflow_run_requested"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "requested",
		Help:      "The total number of workflow runs requested",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_completed"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "completed",
		Help:      "The total number of workflow runs completed",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_success"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "success",
		Help:      "The total number of workflow runs with a 'success' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_cancelled"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "cancelled",
		Help:      "The total number of workflow runs with a 'cancelled' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_action_required"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "action_required",
		Help:      "The total number of workflow runs with a 'action_required' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_timed_out"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "timed_out",
		Help:      "The total number of workflow runs with a 'timed_out' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_failure"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "failure",
		Help:      "The total number of workflow runs with a 'failure' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_neutral"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "neutral",
		Help:      "The total number of workflow runs with a 'neutral' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_skipped"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "skipped",
		Help:      "The total number of workflow runs with a 'skipped' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_startup_failure"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "startup_failure",
		Help:      "The total number of workflow runs with a 'startup_failure' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_run_stale"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "stale",
		Help:      "The total number of workflow runs with a 'stale' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	histograms["github_actions_workflow_run_duration"] = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "duration",
		Help:      "The duration of workflow runs",
		Buckets:   prometheus.LinearBuckets(0, 2, 10),
	},
		[]string{"org", "repo", "workflow"},
	)

	counters["github_actions_workflow_job_queued"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "queued",
		Help:      "The total number of workflow jobs queued",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	counters["github_actions_workflow_job_in_progress"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "in_progress",
		Help:      "The total number of workflow jobs in progress",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	counters["github_actions_workflow_job_completed"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "completed",
		Help:      "The total number of workflow jobs completed",
	},
		[]string{"org", "repo", "workflow", "job"},
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
		reportWorkflowRunEvent(github.WebHookType(req), e)
	case *github.WorkflowJobEvent:
		reportWorkflowJobEvent(github.WebHookType(req), e)
	default:
		// log.Printf("unknown event type %s\n", github.WebHookType(req))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func cacheWorkflowNames(event *github.WorkflowRunEvent) {
	workflowNames[fmt.Sprint(event.GetWorkflowRun().GetID())] = event.GetWorkflow().GetName()
}

func getCachedWorkflowNames(runId string) string {
	return workflowNames[runId]
}

func logWorkflowRunEvent(eventType string, event *github.WorkflowRunEvent) {
	log.Printf("reporting %s event with action %s for %s/%s on %s\n", eventType, event.GetAction(), event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName())
}

func logWorkflowJobEvent(eventType string, event *github.WorkflowJobEvent) {
	log.Printf("reporting %s event with action %s for %s/%s on %s\n", eventType, event.GetAction(), event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflowJob().GetName())
}

func getCounter(eventType string, action string) (*prometheus.CounterVec, bool) {
	metricName := fmt.Sprintf("github_actions_%s_%s", eventType, action)
	metric, found := counters[metricName]

	if !found {
		log.Printf("metric not registered %s\n", metricName)
		return nil, false
	}

	return metric, found
}

func getHistogram(eventType string, name string) (*prometheus.HistogramVec, bool) {
	metricName := fmt.Sprintf("github_actions_%s_%s", eventType, name)
	metric, found := histograms[metricName]

	if !found {
		log.Printf("metric not registered %s\n", metricName)
		return nil, false
	}

	return metric, found
}

func reportWorkflowRunEvent(eventType string, event *github.WorkflowRunEvent) {
	logWorkflowRunEvent(eventType, event)
	actionCounter, found := getCounter(eventType, event.GetAction())

	if !found {
		return
	}
	cacheWorkflowNames(event)
	actionCounter.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName()).Inc()

	if event.GetAction() == "completed" {
		conclusionCounter, found := getCounter(eventType, event.GetWorkflowRun().GetConclusion())
		if found {
			conclusionCounter.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName()).Inc()
		}

		histogram, found := getHistogram(eventType, "duration")
		if found {
			// This is elapse time, not billing time. Billing time is the sum of the time spent in each job.
			start := event.GetWorkflowRun().GetCreatedAt().Time
			end := event.GetWorkflowRun().GetUpdatedAt().Time

			histogram.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName()).Observe(float64(end.Sub(start).Milliseconds()))
		}
	}
}

func reportWorkflowJobEvent(eventType string, event *github.WorkflowJobEvent) {
	logWorkflowJobEvent(eventType, event)
	github_actions_workflow_counter, found := getCounter(eventType, event.GetAction())
	if !found {
		return
	}
	workflowName := getCachedWorkflowNames(fmt.Sprint(event.GetWorkflowJob().GetRunID()))
	if workflowName == "" {
		log.Printf("could not find workflow name in cache for workflow run %d\n", event.GetWorkflowJob().GetRunID())
		return
	}
	github_actions_workflow_counter.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), workflowName, event.GetWorkflowJob().GetName()).Inc()
}
