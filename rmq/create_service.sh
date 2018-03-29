#!/usr/bin/env bash
set -eu
: ${RMQ_SERVICE_NAME:=test-rmq}
cf create-service rabbitmq standard "$RMQ_SERVICE_NAME"
