package main

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricSet struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
	cache      WorkflowNameCache
}

func (ms MetricSet) getCounter(eventType string, action string) (*prometheus.CounterVec, bool) {
	metricName := fmt.Sprintf("github_actions_%s_%s", eventType, action)
	metric, found := ms.counters[metricName]

	if !found {
		log.Printf("metric not registered %s\n", metricName)
		return nil, false
	}

	return metric, found
}

func (ms MetricSet) getHistogram(eventType string, name string) (*prometheus.HistogramVec, bool) {
	metricName := fmt.Sprintf("github_actions_%s_%s", eventType, name)
	metric, found := ms.histograms[metricName]

	if !found {
		log.Printf("metric not registered %s\n", metricName)
		return nil, false
	}

	return metric, found
}
