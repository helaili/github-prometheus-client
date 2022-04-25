package main

import (
	"fmt"
	"log"

	"github.com/google/go-github/v43/github"
)

type LocalCache struct {
	AbstractCache
	workflowNames map[string]string
}

func NewLocalCache(app_id int64, private_key []byte) *LocalCache {
	log.Println("Using the local cache")
	return &LocalCache{
		*NewAbstractCache(app_id, private_key),
		map[string]string{},
	}
}

func (m LocalCache) set(event *github.WorkflowRunEvent) {
	m.workflowNames[fmt.Sprintf("%d-%d", event.GetInstallation().GetID(), event.GetWorkflowRun().GetID())] = event.GetWorkflow().GetName()
}

func (m LocalCache) get(event *github.WorkflowJobEvent) string {
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
