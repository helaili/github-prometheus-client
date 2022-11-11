package main

import (
	"fmt"
	"log"

	"github.com/google/go-github/v43/github"
	"github.com/prometheus/client_golang/prometheus"
)

type WorkflowRunMetrics struct {
	MetricSet
}

func NewWorkflowRunMetrics(registry *prometheus.Registry, cache ICache) *WorkflowRunMetrics {
	m := new(WorkflowRunMetrics)
	m.counters = make(map[string]*prometheus.CounterVec)
	m.histograms = make(map[string]*prometheus.HistogramVec)
	m.gauges = make(map[string]*prometheus.GaugeVec)
	m.cache = cache
	m.registry = registry

	m.intializeCounters()
	m.intializeHistograms()
	m.intializeGauges()

	return m
}

func (m WorkflowRunMetrics) log(eventType string, event *github.WorkflowRunEvent) {
	log.Printf("reporting %s event with action %s for %s/%s on %s for installation %d\n", eventType, event.GetAction(), event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName(), event.GetInstallation().GetID())
}

func (m WorkflowRunMetrics) report(eventType string, event *github.WorkflowRunEvent) {
	m.log(eventType, event)
	actionCounter, found := m.getCounter(eventType, event.GetAction())

	installationID := fmt.Sprint(event.GetInstallation().GetID())

	if found {
		m.cache.set(event)
		actionCounter.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName(), installationID).Inc()
	}

	if event.GetAction() == "completed" {
		conclusionCounter, found := m.getCounter(eventType, event.GetWorkflowRun().GetConclusion())
		if found {
			conclusionCounter.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName(), installationID).Inc()
		}

		histogram, found := m.getHistogram(eventType, "duration")
		if found {
			// This is elapse time, not billing time.
			start := event.GetWorkflowRun().GetCreatedAt().Time
			end := event.GetWorkflowRun().GetUpdatedAt().Time

			histogram.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName(), installationID).Observe(float64(end.Sub(start).Milliseconds()))
		}

		gauge, found := m.getGauge(eventType, "duration")
		if found {
			// This is elapse time, not billing time.
			start := event.GetWorkflowRun().GetCreatedAt().Time
			end := event.GetWorkflowRun().GetUpdatedAt().Time

			gauge.WithLabelValues(event.GetOrg().GetLogin(), event.GetRepo().GetName(), event.GetWorkflow().GetName(), installationID).Set(float64(end.Sub(start).Milliseconds()))
		}
	}
}

func (m WorkflowRunMetrics) intializeHistograms() {
	m.histograms["github_actions_workflow_run_duration"] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "duration",
		Help:      "The duration of workflow runs",
		Buckets:   prometheus.LinearBuckets(0, 2, 10),
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	for histogramName := range m.histograms {
		m.registry.MustRegister(m.histograms[histogramName])
	}
}

func (m WorkflowRunMetrics) intializeGauges() {
	m.gauges["github_actions_workflow_run_duration_gauge"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "duration_gauge",
		Help:      "The duration of workflow runs",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	for gaugeName := range m.gauges {
		m.registry.MustRegister(m.gauges[gaugeName])
	}
}

func (m WorkflowRunMetrics) intializeCounters() {
	m.counters["github_actions_workflow_run_requested"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "requested",
		Help:      "The total number of workflow runs requested",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_completed"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "completed",
		Help:      "The total number of workflow runs completed",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_success"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "success",
		Help:      "The total number of workflow runs with a 'success' conclusion",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_cancelled"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "cancelled",
		Help:      "The total number of workflow runs with a 'cancelled' conclusion",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_action_required"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "action_required",
		Help:      "The total number of workflow runs with a 'action_required' conclusion",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_timed_out"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "timed_out",
		Help:      "The total number of workflow runs with a 'timed_out' conclusion",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_failure"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "failure",
		Help:      "The total number of workflow runs with a 'failure' conclusion",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_neutral"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "neutral",
		Help:      "The total number of workflow runs with a 'neutral' conclusion",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_skipped"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "skipped",
		Help:      "The total number of workflow runs with a 'skipped' conclusion",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_startup_failure"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "startup_failure",
		Help:      "The total number of workflow runs with a 'startup_failure' conclusion",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	m.counters["github_actions_workflow_run_stale"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "github_actions",
		Subsystem: "workflow_run",
		Name:      "stale",
		Help:      "The total number of workflow runs with a 'stale' conclusion",
	},
		[]string{"org", "repo", "workflow", "installation"},
	)

	for counterName := range m.counters {
		m.registry.MustRegister(m.counters[counterName])
	}
}
