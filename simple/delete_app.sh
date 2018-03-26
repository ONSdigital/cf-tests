#!/usr/bin/env bash
set -eu
: ${SIMPLE_APP_NAME:=cf-test-simple}
cf delete "$SIMPLE_APP_NAME" -f
