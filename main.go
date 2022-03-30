package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v43/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	http.HandleFunc("/webhook", webhook)

	fmt.Printf("Listening on address %s\n", *addr)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func webhook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s event on end point %s\n", r.Header.Get("X-GitHub-Event"), r.URL)

	payload, err := github.ValidatePayload(r, []byte(""))
	if err != nil {
		log.Printf("error reading request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	//switch e := event.(type) {
	switch event.(type) {
	case *github.WorkflowRunEvent:
		fmt.Println("It's a workflow run!")
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}

	w.WriteHeader(http.StatusOK)
}
