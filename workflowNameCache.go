package main

import (
	"context"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v43/github"
)

type IWorkflowNameCache interface {
	set(event *github.WorkflowRunEvent)
	get(event *github.WorkflowJobEvent) string
}

type WorkflowNameCache struct {
	transport *ghinstallation.AppsTransport
}

func NewWorkflowNameCache(app_id int64, private_key []byte) *WorkflowNameCache {
	transport, err := ghinstallation.NewAppsTransport(http.DefaultTransport, app_id, private_key)
	if err != nil {
		log.Fatal("Failed to initialize GitHub App transport:", err)
	}

	return &WorkflowNameCache{
		transport,
	}
}

func (m WorkflowNameCache) getWorkflowNameFromGitHub(event *github.WorkflowJobEvent) string {
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
	return worflowRun.GetName()
}
