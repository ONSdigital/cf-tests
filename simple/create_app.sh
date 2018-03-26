#!/usr/bin/env bash
set -eu
: ${SIMPLE_APP_NAME:=cf-test-simple}
cf push "$SIMPLE_APP_NAME"
