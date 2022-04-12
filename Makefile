SHELL := /bin/bash
UNAME := $(shell uname)

.PHONY: build clean run-local test

build: clean
	go build -o github-prometheus-client

run-local:
	export GITHUB_PROMETHEUS_CLIENT_ENV=development && ./github-prometheus-client

clean:
	rm -f github-prometheus-client

test:
	cd test && ./test.sh && cd ..
