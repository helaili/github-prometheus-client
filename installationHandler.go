package main

import (
	"log"

	"github.com/google/go-github/v43/github"
)

type InstallationHandler struct {
}

func NewInstallationHandler() *InstallationHandler {
	m := new(InstallationHandler)
	return m
}

func (m InstallationHandler) created(eventType string, event *github.InstallationEvent) {
	log.Printf("reporting %s event with action %s for  %s for installation %d\n", eventType, event.GetAction(), event.GetInstallation().GetAccount().GetLogin(), event.GetInstallation().GetID())
}
