package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/google/go-github/v43/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var workflowRunMetrics *WorkflowRunMetrics
var workflowJobMetrics *WorkflowJobMetrics
var installationHandler *InstallationHandler
var webhook_secret []byte

func main() {
	initialize()

	port := os.Getenv("PORT")

	// This is the Prometheus endpoint.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	// This is the GitHub Webhook endpoint.
	http.HandleFunc("/webhook", webhook)

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func initialize() {
	env, private_key, secret, app_id := initializeEnv()
	webhook_secret = secret
	initializeWorkflowMetrics(env, private_key, app_id)
	initializeInstallationHandler()
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

func initializeWorkflowMetrics(env, private_key string, app_id int64) {
	var workflowNames IWorkflowNameCache
	if env == "development" {
		workflowNames = NewWorkflowNameLocalCache(app_id, []byte(private_key))
	} else {
		workflowNames = NewWorkflowNameRedisCache(app_id, []byte(private_key))
	}
	workflowRunMetrics = NewWorkflowRunMetrics(workflowNames)
	workflowJobMetrics = NewWorkflowJobMetrics(workflowNames)
}

func initializeInstallationHandler() {
	installationHandler = NewInstallationHandler()
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
		workflowRunMetrics.report(github.WebHookType(req), e)
	case *github.WorkflowJobEvent:
		workflowJobMetrics.report(github.WebHookType(req), e)
	case *github.InstallationEvent:
		installationHandler.created(github.WebHookType(req), e)
	default:
		// log.Printf("unknown event type %s\n", github.WebHookType(req))
		return
	}

	w.WriteHeader(http.StatusOK)
}
