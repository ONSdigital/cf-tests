#!/usr/bin/env bash
set -eu
: ${RMQ_SERVICE_NAME:=test-rmq}
cf service "$RMQ_SERVICE_NAME"
