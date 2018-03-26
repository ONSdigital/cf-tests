#!/usr/bin/env bash
set -eu
: ${ELASTICACHE_SERVICE_NAME:=test-elasticache}
cf create-service elasticache-broker small "$ELASTICACHE_SERVICE_NAME"
