package main

import (
	"log"

	"github.com/google/go-github/v43/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type WorkflowJobMetrics struct {
	MetricSet
}

func NewWorkflowJobMetrics(cache WorkflowNameCache) *WorkflowJobMetrics {
	m := new(WorkflowJobMetrics)
	m.counters = make(map[string]*prometheus.CounterVec)
	m.histograms = make(map[string]*prometheus.HistogramVec)
	m.cache = cache

	m.intializeCounters()
	m.intializeHistograms()

	return m
}

func (m WorkflowJobMetrics) log(eventType string, event *github.WorkflowJobEvent) {
	log.Printf("reporting %s event with action %s for %s/%s on %s\n", eventType, event.GetAction(), event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflowJob().GetName())
}

func (m WorkflowJobMetrics) report(eventType string, event *github.WorkflowJobEvent) {
	m.log(eventType, event)

	workflowName := m.cache.get(event)
	if workflowName == "" {
		log.Printf("could not find workflow name in cache for workflow run %d\n", event.GetWorkflowJob().GetRunID())
		workflowName = "unknown"
	}

	actionCounter, found := m.getCounter(eventType, event.GetAction())
	if found {
		actionCounter.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), workflowName, event.GetWorkflowJob().GetName()).Inc()
	}

	if event.GetAction() == "completed" {
		conclusionCounter, found := m.getCounter(eventType, event.GetWorkflowJob().GetConclusion())
		if found {
			conclusionCounter.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), workflowName, event.GetWorkflowJob().GetName()).Inc()
		}

		histogram, found := m.getHistogram(eventType, "duration")
		if found {
			// This is the billing time.
			start := event.GetWorkflowJob().GetStartedAt().Time
			end := event.GetWorkflowJob().GetCompletedAt().Time

			histogram.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), workflowName, event.GetWorkflowJob().GetName()).Observe(float64(end.Sub(start).Milliseconds()))
		}
	}
}

func (m WorkflowJobMetrics) intializeHistograms() {
	m.histograms["github_actions_workflow_job_duration"] = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "duration",
		Help:      "The duration of workflow job, equivalent to the billing time",
		Buckets:   prometheus.LinearBuckets(0, 2, 10),
	},
		[]string{"org", "repo", "workflow", "job"},
	)
}

func (m WorkflowJobMetrics) intializeCounters() {
	m.counters["github_actions_workflow_job_queued"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "queued",
		Help:      "The total number of workflow jobs queued",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_in_progress"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "in_progress",
		Help:      "The total number of workflow jobs in progress",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_completed"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "completed",
		Help:      "The total number of workflow jobs completed",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_success"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "success",
		Help:      "The total number of workflow runs with a 'success' conclusion",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_cancelled"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "cancelled",
		Help:      "The total number of workflow runs with a 'cancelled' conclusion",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_action_required"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "action_required",
		Help:      "The total number of workflow runs with a 'action_required' conclusion",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_timed_out"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "timed_out",
		Help:      "The total number of workflow runs with a 'timed_out' conclusion",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_failure"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "failure",
		Help:      "The total number of workflow runs with a 'failure' conclusion",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_neutral"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "neutral",
		Help:      "The total number of workflow runs with a 'neutral' conclusion",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_skipped"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "skipped",
		Help:      "The total number of workflow runs with a 'skipped' conclusion",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_startup_failure"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "startup_failure",
		Help:      "The total number of workflow runs with a 'startup_failure' conclusion",
	},
		[]string{"org", "repo", "workflow", "job"},
	)

	m.counters["github_actions_workflow_job_stale"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_job",
		Name:      "stale",
		Help:      "The total number of workflow runs with a 'stale' conclusion",
	},
		[]string{"org", "repo", "workflow", "job"},
	)
}
