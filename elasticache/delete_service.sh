#!/usr/bin/env bash
set -eu
: ${ELASTICACHE_SERVICE_NAME:=test-elasticache}
# FIXME: need to unbind the application first
cf delete-service "$ELASTICACHE_SERVICE_NAME" -f
