Run locally: 

```
ngrok http -subdomain=github-prometheus-client 8080
``` 


### Dev mode on local machine
This mode is fully offline. To use it, set the environment variable `GITHUB_PROMETHEUS_CLIENT_ENV` to `development` and then follow the instructions below. 

- Install Prometheus and set the environment variable `PROMETHEUS_PATH`
- Build with `make build`
- Run with `make run`
- In a second shell window, launch the local Prometheus server: `cd test && ./prometheus.sh`
- In a third shell window, launch tests by running `cd test && ./test.sh` 
- Open a browser on `http://localhost:9090/` and add some of the metrics to a graph 
