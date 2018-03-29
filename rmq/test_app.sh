#!/usr/bin/env bash
set -eux
: ${RMQ_ENDPOINT:=https://cf-test-rmq.apps.devtest.onsclofo.uk}
curl -s -k "$RMQ_ENDPOINT" | grep "RMQ service is OK"
