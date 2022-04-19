package main

import (
	"fmt"

	"github.com/google/go-github/v43/github"
)

type WorkflowNameLocalCache struct {
	WorkflowNameCache
	workflowNames map[string]string
}

func NewWorkflowNameLocalCache(app_id int64, private_key []byte) *WorkflowNameLocalCache {
	return &WorkflowNameLocalCache{
		*NewWorkflowNameCache(app_id, private_key),
		map[string]string{},
	}
}

func (m WorkflowNameLocalCache) set(event *github.WorkflowRunEvent) {
	m.workflowNames[fmt.Sprintf("%d-%d", event.GetInstallation().GetID(), event.GetWorkflowRun().GetID())] = event.GetWorkflow().GetName()
}

func (m WorkflowNameLocalCache) get(event *github.WorkflowJobEvent) string {
	runId := fmt.Sprintf("%d-%d", event.GetInstallation().GetID(), event.GetWorkflowJob().GetRunID())
	worfklowName, ok := m.workflowNames[runId]

	if ok {
		return worfklowName
	} else {
		workflowName := m.getWorkflowNameFromGitHub(event)
		m.workflowNames[runId] = workflowName
		return workflowName
	}
}
