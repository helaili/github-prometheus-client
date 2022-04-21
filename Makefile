SHELL := /bin/sh
UNAME := $(shell uname)

.PHONY: build clean run run-dev test e2e-test docker-build

build: clean
	go build -o github-prometheus-client

run-dev:
	export GITHUB_PROMETHEUS_CLIENT_ENV=development && ./github-prometheus-client

run:
	./github-prometheus-client

clean:
	rm -f github-prometheus-client

test:
	go test

e2e-test:
	cd test && ./test.sh && cd ..

docker-build:
	docker build -t github-prometheus-client .

