#!/usr/bin/env bash
set -eux
: ${SIMPLE_ENDPOINT:=https://cf-test-simple.apps.devtest.onsclofo.uk}
curl -s -k "$SIMPLE_ENDPOINT" | grep "Test app is OK"
