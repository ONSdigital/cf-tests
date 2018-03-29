#!/usr/bin/env bash
set -eu
: ${RMQ_APP_NAME:=cf-test-rmq}
cf push "$RMQ_APP_NAME"
