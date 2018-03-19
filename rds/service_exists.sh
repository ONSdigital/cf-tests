#!/usr/bin/env bash
set -eu
: ${RDS_SERVICE_NAME:=test-psql}
cf service "$RDS_SERVICE_NAME"
