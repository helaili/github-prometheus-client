package main

import (
	"fmt"

	"github.com/google/go-github/v43/github"
)

type WorkflowNameCache interface {
	set(event *github.WorkflowRunEvent)
	get(runId string) string
}

type WorkflowNameCacheImpl struct {
	// TODO: use a distributed cache so that we can expire and make the data available to other instances
	workflowNames map[string]string
}

func NewWorkflowNameCacheImpl() *WorkflowNameCacheImpl {
	return &WorkflowNameCacheImpl{
		map[string]string{},
	}
}

func (m WorkflowNameCacheImpl) set(event *github.WorkflowRunEvent) {
	m.workflowNames[fmt.Sprint(event.GetWorkflowRun().GetID())] = event.GetWorkflow().GetName()
}

func (m WorkflowNameCacheImpl) get(runId string) string {
	return m.workflowNames[runId]
}
