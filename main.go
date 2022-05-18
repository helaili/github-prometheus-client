package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/joho/godotenv"

	"github.com/google/go-github/v43/github"
)

// Each installation has a distinct handler so we can isolate
// collectors in separate Prometheus registries and listen on
// a dedicated URL per Prometheus scraper.
var installationHandlers map[string]*InstallationHandler
var webhook_secret []byte
var cache ICache

func main() {
	env, private_key, secret, app_id := initializeEnv()
	webhook_secret = secret
	port := os.Getenv("PORT")

	installationHandlers = make(map[string]*InstallationHandler)

	if env == "development" || env == "dev" {
		cache = NewLocalCache(app_id, []byte(private_key))
		initializeDummyInstallationHandlers()
	} else {
		redisAddress := os.Getenv("REDIS_ADDRESS")
		redisPassword := os.Getenv("REDIS_PASSWORD")
		cache = NewRedisCache(redisAddress, redisPassword, app_id, []byte(private_key))
		initializeInstallationHandlers(app_id, []byte(private_key))
	}

	http.HandleFunc("/ping", ping)

	// This is the GitHub Webhook endpoint.
	http.HandleFunc("/webhook", webhook)

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func initializeEnv() (env string, private_key string, webhook_secret []byte, app_id int64) {
	env = os.Getenv("GITHUB_PROMETHEUS_CLIENT_ENV")
	if env == "" || env == "dev" {
		env = "development"
	}

	log.Printf("Starting in %s mode\n", env)

	godotenv.Load(".env." + env)
	godotenv.Load()

	if env == "development" {
		return env, "dummy", []byte("dummy"), 0
	} else {
		private_key = os.Getenv("PRIVATE_KEY")
		// Private key a one line string. For some reasons, '\n' are not interpreted correctly.
		// We there for provide a string in the environment where '\n's are replaced with "^"s.
		// We now need to put these \n in.
		private_key = strings.Replace(private_key, "^", "\n", -1)

		app_id, err := strconv.ParseInt(os.Getenv("APP_ID"), 10, 36)
		if err != nil {
			log.Fatal("Wrong format for APP_ID")
		}
		webhook_secret = []byte(os.Getenv("WEBHOOK_SECRET"))
		return env, private_key, webhook_secret, app_id
	}
}

/*
 * Create a hanlder for each existing installation.
 */
func initializeInstallationHandlers(app_id int64, private_key []byte) {
	transport, err := ghinstallation.NewAppsTransport(http.DefaultTransport, app_id, private_key)
	if err != nil {
		log.Fatal("Failed to initialize GitHub App transport:", err)
	}

	client := github.NewClient(&http.Client{Transport: transport})

	listOptions := github.ListOptions{PerPage: 100}
	listOptions.Page = 1

	for listOptions.Page != 0 {
		installations, res, err := client.Apps.ListInstallations(context.Background(), &listOptions)
		if err != nil {
			log.Fatal("Failed to retrieve App installations:", err)
		}

		for _, installation := range installations {
			log.Printf("Initializing installation %d\n", installation.GetID())
			installationHandlers[fmt.Sprintf("%d", installation.GetID())] = NewInstallationHandler(installation.GetID(), cache)
		}

		listOptions.Page = res.NextPage
	}
}

/*
 * Initialize static handlers matching test datat for offline testing
 */
func initializeDummyInstallationHandlers() {
	installationHandlers["24886277"] = NewInstallationHandler(24886277, cache)
	installationHandlers["25140335"] = NewInstallationHandler(25140335, cache)
}

func getInstallationHandler(installation_id int64) *InstallationHandler {
	handler := installationHandlers[fmt.Sprintf("%d", installation_id)]
	if handler == nil {
		log.Printf("No handler for installation %d\n", installation_id)
	}
	return handler
}

func ping(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Ok"))
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
		handler := getInstallationHandler(e.GetInstallation().GetID())
		if handler != nil {
			handler.workflowRunMetrics.report(github.WebHookType(req), e)
		}
	case *github.WorkflowJobEvent:
		handler := getInstallationHandler(e.GetInstallation().GetID())
		if handler != nil {
			handler.workflowJobMetrics.report(github.WebHookType(req), e)
		}
	case *github.InstallationEvent:
		installationHandlers[fmt.Sprintf("%d", e.GetInstallation().GetID())] = NewInstallationHandler(e.GetInstallation().GetID(), cache)
	default:
		// log.Printf("unknown event type %s\n", github.WebHookType(req))
		return
	}

	w.WriteHeader(http.StatusOK)
}
