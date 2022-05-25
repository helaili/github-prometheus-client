package main

import (
	"fmt"
	"log"

	"github.com/google/go-github/v43/github"
)

type OfflineCache struct {
	AbstractCache
	workflowNames map[string]string
}

func NewOfflineCache(app_id int64) *OfflineCache {
	log.Println("Using the offline cache")
	return &OfflineCache{
		AbstractCache{},
		map[string]string{},
	}
}

func (m OfflineCache) set(event *github.WorkflowRunEvent) {
	m.workflowNames[fmt.Sprintf("%d-%d", event.GetInstallation().GetID(), event.GetWorkflowRun().GetID())] = event.GetWorkflow().GetName()
}

func (m OfflineCache) get(event *github.WorkflowJobEvent) string {
	runId := fmt.Sprintf("%d-%d", event.GetInstallation().GetID(), event.GetWorkflowJob().GetRunID())
	worfklowName := m.workflowNames[runId]
	return worfklowName
}
