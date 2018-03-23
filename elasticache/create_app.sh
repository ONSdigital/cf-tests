#!/usr/bin/env bash
set -eu
: ${ELASTICACHE_APP_NAME:=cf-test-elasticache}
cf push "$ELASTICACHE_APP_NAME"
