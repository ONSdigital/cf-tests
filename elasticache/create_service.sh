#!/usr/bin/env bash
set -eu
: ${ELASTICACHE_SERVICE_NAME:=test-elasticache}
cf create-service elasticache shared-elasticache "$ELASTICACHE_SERVICE_NAME"
