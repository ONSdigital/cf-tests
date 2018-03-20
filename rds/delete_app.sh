#!/usr/bin/env bash
set -eu
: ${RDS_APP_NAME:=cf-test-rds}
cf delete "$RDS_APP_NAME" -f
