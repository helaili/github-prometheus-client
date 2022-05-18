########################################################################
# Run a local instance of Prometheus with the right config for testing #
########################################################################
#!/bin/sh
set -e

$PROMETHEUS_PATH/prometheus --config.file=prometheus/prometheus-24886277.yml --web.listen-address="0.0.0.0:9091" --storage.tsdb.path="prometheus/data-24886278/"