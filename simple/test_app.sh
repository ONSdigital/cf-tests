#!/usr/bin/env bash
set -eu
: ${SIMPLE_ENDPOINT:=https://cf-test-simple.apps.devtest.onsclofo.uk}
curl -k "$SIMPLE_ENDPOINT"
