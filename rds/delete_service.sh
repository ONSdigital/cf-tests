#!/usr/bin/env bash
set -eu
: ${RDS_SERVICE_NAME:=test-psql}
# FIXME: need to unbind the application first
cf delete-service "$RDS_SERVICE_NAME" -f
