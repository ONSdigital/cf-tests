#!/usr/bin/env bash
set -eu
: ${ELASTICACHE_ENDPOINT:=https://cf-test-elasticache.apps.devtest.onsclofo.uk}
curl -k "$ELASTICACHE_ENDPOINT"
