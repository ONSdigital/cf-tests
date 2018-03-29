#!/usr/bin/env bash
set -eu
: ${RMQ_APP_NAME:=cf-test-rmq}
cf delete "$RMQ_APP_NAME" -f
