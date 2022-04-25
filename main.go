package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/google/go-github/v43/github"
)

// Each installation has a distinct handler so we can isolate
// collectors in separate Prometheus registries and listen on
// a dedicated URL per Prometheus scraper.
var installationHandlers map[string]*InstallationHandler
var webhook_secret []byte

func main() {
	env, private_key, secret, app_id := initializeEnv()
	webhook_secret = secret
	port := os.Getenv("PORT")

	var cache IWorkflowNameCache
	if env == "development" {
		cache = NewWorkflowNameLocalCache(app_id, []byte(private_key))
	} else {
		cache = NewWorkflowNameRedisCache(app_id, []byte(private_key))
	}

	installationHandlers = make(map[string]*InstallationHandler)
	initializeInstallationHandlers(cache)

	// This is the GitHub Webhook endpoint.
	http.HandleFunc("/webhook", webhook)

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func initializeEnv() (env string, private_key string, webhook_secret []byte, app_id int64) {
	env = os.Getenv("GITHUB_PROMETHEUS_CLIENT_ENV")
	if env == "" {
		env = "development"
	}

	log.Printf("Starting in %s mode\n", env)

	godotenv.Load(".env." + env)
	godotenv.Load()

	private_key = os.Getenv("PRIVATE_KEY")
	app_id, err := strconv.ParseInt(os.Getenv("APP_ID"), 10, 36)
	if err != nil {
		log.Fatal("Wrong format for APP_ID")
	}
	webhook_secret = []byte(os.Getenv("WEBHOOK_SECRET"))

	return env, private_key, webhook_secret, app_id
}

func initializeInstallationHandlers(cache IWorkflowNameCache) {
	installationHandlers[fmt.Sprintf("%d", 24886277)] = NewInstallationHandler(24886277, cache)
	installationHandlers[fmt.Sprintf("%d", 24886278)] = NewInstallationHandler(24886278, cache)
}

func getInstallationHandler(installation_id int64) *InstallationHandler {
	handler := installationHandlers[fmt.Sprintf("%d", installation_id)]
	if handler == nil {
		log.Printf("No handler for installation %d\n", installation_id)
	}
	return handler
}

/*
 * Receive a webhook from GitHub and report the event to Prometheus.
 */
func webhook(w http.ResponseWriter, req *http.Request) {
	log.Printf("Received %s event on end point %s\n", req.Header.Get("X-GitHub-Event"), req.URL)

	if webhook_secret == nil {
		webhook_secret = []byte{}
	}

	payload, err := github.ValidatePayload(req, webhook_secret)
	if err != nil {
		log.Printf("error reading request body: err=%s\n", err)
		return
	}
	defer req.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(req), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.WorkflowRunEvent:
		getInstallationHandler(e.GetInstallation().GetID()).workflowRunMetrics.report(github.WebHookType(req), e)
	case *github.WorkflowJobEvent:
		getInstallationHandler(e.GetInstallation().GetID()).workflowJobMetrics.report(github.WebHookType(req), e)
	case *github.InstallationEvent:
		//installationHandler.created(github.WebHookType(req), e)
	default:
		// log.Printf("unknown event type %s\n", github.WebHookType(req))
		return
	}

	w.WriteHeader(http.StatusOK)
}
