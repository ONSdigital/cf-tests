#!/usr/bin/env bash
set -eux
: ${ELASTICACHE_ENDPOINT:=https://cf-test-elasticache.apps.devtest.onsclofo.uk}
curl -s -k "$ELASTICACHE_ENDPOINT" | grep "Elasticache service is OK"
