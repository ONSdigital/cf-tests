#!/usr/bin/env bash
set -eu
: ${RMQ_SERVICE_NAME:=test-rmq}
# FIXME: need to unbind the application first
cf delete-service "$RMQ_SERVICE_NAME" -f
