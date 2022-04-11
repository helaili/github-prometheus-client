package main

import (
	"log"

	"github.com/google/go-github/v43/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type WorkflowRunMetrics struct {
	MetricSet
}

func NewWorkflowRunMetrics(cache WorkflowNameCache) *WorkflowRunMetrics {
	m := new(WorkflowRunMetrics)
	m.counters = make(map[string]*prometheus.CounterVec)
	m.histograms = make(map[string]*prometheus.HistogramVec)
	m.cache = cache

	m.intializeCounters()
	m.intializeHistograms()

	return m
}

func (m WorkflowRunMetrics) log(eventType string, event *github.WorkflowRunEvent) {
	log.Printf("reporting %s event with action %s for %s/%s on %s\n", eventType, event.GetAction(), event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName())
}

func (m WorkflowRunMetrics) report(eventType string, event *github.WorkflowRunEvent) {
	m.log(eventType, event)
	actionCounter, found := m.getCounter(eventType, event.GetAction())

	if found {
		m.cache.set(event)
		actionCounter.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName()).Inc()
	}

	if event.GetAction() == "completed" {
		conclusionCounter, found := m.getCounter(eventType, event.GetWorkflowRun().GetConclusion())
		if found {
			conclusionCounter.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName()).Inc()
		}

		histogram, found := m.getHistogram(eventType, "duration")
		if found {
			// This is elapse time, not billing time.
			start := event.GetWorkflowRun().GetCreatedAt().Time
			end := event.GetWorkflowRun().GetUpdatedAt().Time

			histogram.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName()).Observe(float64(end.Sub(start).Milliseconds()))
		}
	}
}

func (m WorkflowRunMetrics) intializeHistograms() {
	m.histograms["github_actions_workflow_run_duration"] = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "duration",
		Help:      "The duration of workflow runs",
		Buckets:   prometheus.LinearBuckets(0, 2, 10),
	},
		[]string{"org", "repo", "workflow"},
	)
}

func (m WorkflowRunMetrics) intializeCounters() {
	m.counters["github_actions_workflow_run_requested"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "requested",
		Help:      "The total number of workflow runs requested",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_completed"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "completed",
		Help:      "The total number of workflow runs completed",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_success"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "success",
		Help:      "The total number of workflow runs with a 'success' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_cancelled"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "cancelled",
		Help:      "The total number of workflow runs with a 'cancelled' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_action_required"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "action_required",
		Help:      "The total number of workflow runs with a 'action_required' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_timed_out"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "timed_out",
		Help:      "The total number of workflow runs with a 'timed_out' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_failure"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "failure",
		Help:      "The total number of workflow runs with a 'failure' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_neutral"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "neutral",
		Help:      "The total number of workflow runs with a 'neutral' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_skipped"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "skipped",
		Help:      "The total number of workflow runs with a 'skipped' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_startup_failure"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "startup_failure",
		Help:      "The total number of workflow runs with a 'startup_failure' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)

	m.counters["github_actions_workflow_run_stale"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "stale",
		Help:      "The total number of workflow runs with a 'stale' conclusion",
	},
		[]string{"org", "repo", "workflow"},
	)
}
