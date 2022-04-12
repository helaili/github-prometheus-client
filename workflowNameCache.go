package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v43/github"
)

type WorkflowNameCache interface {
	set(event *github.WorkflowRunEvent)
	get(event *github.WorkflowJobEvent) string
}

type WorkflowNameCacheImpl struct {
	// TODO: use a distributed cache so that we can expire and make the data available to other instances
	workflowNames map[string]string
	transport     *ghinstallation.AppsTransport
}

func NewWorkflowNameCacheImpl(app_id int64, private_key []byte) *WorkflowNameCacheImpl {
	transport, err := ghinstallation.NewAppsTransport(http.DefaultTransport, app_id, private_key)
	if err != nil {
		log.Fatal("Failed to initialize GitHub App transport:", err)
	}

	return &WorkflowNameCacheImpl{
		map[string]string{},
		transport,
	}
}

func (m WorkflowNameCacheImpl) set(event *github.WorkflowRunEvent) {
	m.workflowNames[fmt.Sprint(event.GetWorkflowRun().GetID())] = event.GetWorkflow().GetName()
}

func (m WorkflowNameCacheImpl) get(event *github.WorkflowJobEvent) string {
	runId := fmt.Sprint(event.GetWorkflowJob().GetRunID())
	worfklowName, ok := m.workflowNames[runId]

	if ok {
		return worfklowName
	} else {
		// The workflow name has not been cached. Let's fetch it from GitHub.
		installationID := event.GetInstallation().GetID()
		if installationID == 0 {
			log.Printf("Failed to retrieve installation ID")
			return ""
		}

		installationTransport := ghinstallation.NewFromAppsTransport(m.transport, installationID)
		client := github.NewClient(&http.Client{Transport: installationTransport})
		worflowRun, _, err := client.Actions.GetWorkflowRunByID(context.Background(), event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), event.GetWorkflowJob().GetRunID())
		if err != nil {
			log.Printf("Failed to retrieve workflow run for %s/%s with id %d: %s", event.GetSender().GetLogin(), event.GetRepo().GetName(), event.GetWorkflowJob().GetRunID(), err)
			return ""
		}

		m.workflowNames[runId] = worflowRun.GetName()

		return worflowRun.GetName()
	}
}
