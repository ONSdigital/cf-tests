#!/usr/bin/env bash
set -eu
: ${ELASTICACHE_SERVICE_NAME:=test-elasticache}
cf service "$ELASTICACHE_SERVICE_NAME"
