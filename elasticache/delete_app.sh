#!/usr/bin/env bash
set -eu
: ${ELASTICACHE_APP_NAME:=cf-test-elasticache}
cf delete "$ELASTICACHE_APP_NAME" -f
